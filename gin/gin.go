package gin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sp/config"
	"sp/db"
	"sp/ethclient"
	"sp/types"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

// rtCache 实时链上数据快照；刷新时在新对象上构建、构建完成后原子替换，
// 读侧始终拿到完整的旧快照，刷新过程不阻断查询。
type rtCache struct {
	points  []types.SoulPoints // points > 0, 按 points 降序
	byAddr  map[common.Address]types.SoulPoints
	allJSON []byte // 全量响应预序列化，请求直接写出
}

const rtSwitchTime = 1786377600 // 2026/8/11 00:00 UTC+8，之后 / 切换为链上实时口径

func levelOf(points uint64) string {
	switch {
	case points >= 710:
		return "Platinum"
	case points >= 480:
		return "Gold"
	case points >= 48:
		return "Entry"
	case points >= 1:
		return "Member"
	default:
		return ""
	}
}

// handlerLevelsFunc 静态等级表，不含用户数据
func handlerLevelsFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "", "data": []gin.H{
		{"level": "Member", "min": 1, "max": 47},
		{"level": "Entry", "min": 48, "max": 479},
		{"level": "Gold", "min": 480, "max": 709},
		{"level": "Platinum", "min": 710, "max": nil},
	}})
}

var rtData atomic.Pointer[rtCache]

const rtRefreshInterval = 15 * time.Minute

func buildRTCache(ctx context.Context) (*rtCache, error) {
	cctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	addrs, err := ethclient.GetAllMembers(cctx, nil)
	if err != nil {
		return nil, err
	}
	details, err := ethclient.GetAllSp(cctx, addrs, nil)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	cache := &rtCache{
		byAddr: make(map[common.Address]types.SoulPoints, len(details)),
	}
	for _, d := range details {
		sp := types.SoulPoints{
			Address:       d.Addr.String(),
			Count:         1,
			Points:        d.Sum().Uint64(),
			KogePoint:     d.KogePoint.Uint64(),
			StakePoint:    d.StakePoint.Uint64(),
			NftPoint:      d.NftPoint.Uint64(),
			BscStakePoint: d.BscStakePoint.Uint64(),
			GovBNBPoint:   d.GovBNBPoint.Uint64(),
			CreatedAt:     now.Unix(),
		}
		sp.Level = levelOf(sp.Points)
		cache.byAddr[d.Addr] = sp
		if sp.Points > 0 {
			cache.points = append(cache.points, sp)
		}
	}
	sort.Slice(cache.points, func(i, j int) bool {
		return cache.points[i].Points > cache.points[j].Points
	})

	allJSON, err := json.Marshal(gin.H{"code": http.StatusOK, "message": "", "data": cache.points})
	if err != nil {
		return nil, err
	}
	cache.allJSON = allJSON
	return cache, nil
}

// refreshRTCache singleflight 合并并发刷新，完成后原子替换快照。
func refreshRTCache(ctx context.Context) (*rtCache, error) {
	v, err, _ := sfGroup.Do("rt-refresh", func() (any, error) {
		return buildRTCache(ctx)
	})
	if err != nil {
		return nil, err
	}
	cache := v.(*rtCache)
	rtData.Store(cache)
	return cache, nil
}

func rtRefresher(ctx context.Context) {
	// 切换前无需刷新；提前一个周期预热，保证切换瞬间已有快照
	if wait := time.Until(time.Unix(rtSwitchTime, 0).Add(-rtRefreshInterval)); wait > 0 {
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
	}
	ticker := time.NewTicker(rtRefreshInterval)
	defer ticker.Stop()
	for {
		if _, err := refreshRTCache(ctx); err != nil {
			log.Printf("refresh rt cache error: %v", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// dbHandler 数据库口径 handler（带页面缓存），Run 中初始化
var dbHandler gin.HandlerFunc

func handlerIndexFunc(c *gin.Context) {
	if time.Now().Unix() < rtSwitchTime {
		dbHandler(c)
		return
	}

	var query types.Query
	err := c.BindQuery(&query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid query", "data": []types.SoulPoints{}})
		return
	}

	cache := rtData.Load()
	if cache == nil {
		// 首次未就绪：同步构建一次，并发请求由 singleflight 合并等待
		cache, err = refreshRTCache(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "chain query error", "data": []types.SoulPoints{}})
			return
		}
	}

	if query.Address == "" {
		c.Data(http.StatusOK, "application/json; charset=utf-8", cache.allJSON)
		return
	}

	address := common.HexToAddress(query.Address)
	if !strings.EqualFold(address.Hex(), query.Address) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid address", "data": []types.SoulPoints{}})
		return
	}

	sp, ok := cache.byAddr[address]
	if !ok || sp.Points == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "no points", "data": []types.SoulPoints{}})
		return
	}

	if !query.Detail {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "", "data": []types.SoulPoints{sp}})
		return
	}

	// points/level 实时，detail 明细来自数据库明细缓存
	dbDetail, err := db.DetailOf(address.Hex())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "mysql error", "data": []types.SoulPoints{}})
		return
	}
	details := types.SoulPointsWithDetail{SoulPoints: sp, Detail: dbDetail}
	details.Count = uint64(len(dbDetail))
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "", "data": details})
}

func handlerFunc(c *gin.Context) {
	var query types.Query
	err := c.BindQuery(&query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid query", "data": []types.SoulPoints{}})
		return
	}

	var (
		points        []types.SoulPoints                             // mysql scan result
		res           any                                            // response
		dbQuery       = db.Server.Debug().Model(&types.SoulPoints{}) // mysql query
		errCode       = http.StatusInternalServerError               // response code
		address       = common.HexToAddress(query.Address)
		st            []types.SnapTime // mysql scan result
		snapshotCount uint64
	)

	tx := db.Server.Order("created DESC").Limit(48).Find(&st)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(errCode, gin.H{"code": http.StatusInternalServerError, "message": "mysql error", "data": []types.SoulPoints{}})
		return
	}

	if tx.RowsAffected == 48 && time.Now().Unix() > 1727740800 {
		snapshotCount = 48
	} else if tx.RowsAffected >= 7 {
		snapshotCount = 7
	} else if tx.RowsAffected >= 2 {
		snapshotCount = 2
	} else {
		// no snapshot
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "no points", "data": []types.SoulPoints{}})
		return
	}

	sumQueryStr := fmt.Sprintf("user_id, users.address AS address, SUM(points) DIV %d AS points, COUNT(user_id) AS `count`, created", snapshotCount)

	if query.Address == "" {
		dbQuery = dbQuery.Select(sumQueryStr).Joins("RIGHT JOIN users ON user_id = users.id").Where("created BETWEEN ? AND ? AND points > 0", st[snapshotCount-1].CreatedAt, st[0].CreatedAt).Group("user_id").Order("points DESC")
		errCode = http.StatusOK
	} else {
		if !strings.EqualFold(address.Hex(), query.Address) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid address", "data": []types.SoulPoints{}})
			return
		}
		if query.Detail {
			dbQuery = dbQuery.Select("user_id, users.address AS address, points, koge_point, stake_point, nft_point, bsc_stake_point, gov_bnb_point, created").Order("created DESC")
		} else {
			dbQuery = dbQuery.Select(sumQueryStr).Group("user_id")
		}
		dbQuery = dbQuery.Joins("RIGHT JOIN users ON user_id = users.id").Where("users.address = ? AND created BETWEEN ? AND ?", address.Hex(), st[snapshotCount-1].CreatedAt, st[0].CreatedAt)
	}

	tx = dbQuery.Find(&points)

	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(errCode, gin.H{"code": http.StatusInternalServerError, "message": "mysql error", "data": []types.SoulPoints{}})
		return
	}

	if tx.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "no points", "data": []types.SoulPoints{}})
		return
	}

	for i := range points {
		points[i].Level = levelOf(points[i].Points)
	}

	if query.Address != "" && query.Detail {
		var details = types.SoulPointsWithDetail{}
		details.Address = points[0].Address
		for _, point := range points {
			details.Points += point.Points
			details.Detail = append(details.Detail, types.Detail{
				SnapTime: time.Unix(point.CreatedAt, 0).Format("2006/01/02"),
				Points:   point.Points,
				Koge:     point.KogePoint,
				Stake:    point.StakePoint,
				Nft:      point.NftPoint,
				BscStake: point.BscStakePoint + point.GovBNBPoint,
			})
		}
		details.Count = uint64(len(points))
		details.Points /= snapshotCount
		details.Level = levelOf(details.Points)
		res = details
	} else {
		res = points
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "", "data": res})
}

var (
	srv     *http.Server
	sfGroup singleflight.Group
)

type cachedResp struct {
	Status      int
	ContentType string
	Body        []byte
}

type captureWriter struct {
	gin.ResponseWriter
	body   bytes.Buffer
	status int
}

func (w *captureWriter) Write(b []byte) (int, error)       { return w.body.Write(b) }
func (w *captureWriter) WriteString(s string) (int, error) { return w.body.WriteString(s) }
func (w *captureWriter) WriteHeader(code int)              { w.status = code }
func (w *captureWriter) WriteHeaderNow()                   {}
func (w *captureWriter) Written() bool                     { return w.body.Len() > 0 || w.status != 0 }
func (w *captureWriter) Size() int                         { return w.body.Len() }
func (w *captureWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

// sfCachePage 用 singleflight 防止缓存过期时的击穿：同一 key 并发请求只有一个执行
// 真正的 handler，其余等待并共享结果，之后写入缓存供后续请求直接命中。
func sfCachePage(store persistence.CacheStore, expire time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.URL.RequestURI()

		// 快路径：缓存命中直接返回
		var cached cachedResp
		if err := store.Get(key, &cached); err == nil {
			c.Data(cached.Status, cached.ContentType, cached.Body)
			return
		}

		// 慢路径：singleflight 保证同一 key 只有一个 goroutine 执行 handler
		v, _, _ := sfGroup.Do(key, func() (any, error) {
			// double-check：可能已被另一个 goroutine 填充
			var cached2 cachedResp
			if err := store.Get(key, &cached2); err == nil {
				return &cached2, nil
			}

			origWriter := c.Writer
			cw := &captureWriter{ResponseWriter: origWriter}
			c.Writer = cw
			handle(c)
			c.Writer = origWriter

			resp := &cachedResp{
				Status:      cw.Status(),
				ContentType: origWriter.Header().Get("Content-Type"),
				Body:        cw.body.Bytes(),
			}
			// 只缓存正常响应，避免将错误结果长期缓存
			if resp.Status == http.StatusOK {
				_ = store.Set(key, *resp, expire)
			}
			return resp, nil
		})

		if resp, ok := v.(*cachedResp); ok {
			c.Data(resp.Status, resp.ContentType, resp.Body)
		}
	}
}

func Run(pctx context.Context) chan struct{} {
	done := make(chan struct{})
	r := gin.Default()
	r.TrustedPlatform = gin.PlatformCloudflare

	store := persistence.NewInMemoryStore(time.Second * 10)

	r.Use(addCors(), checkHost)

	dbHandler = sfCachePage(store, time.Hour, handlerFunc)

	r.GET("/", handlerIndexFunc)

	r.GET("/levels", handlerLevelsFunc)

	// 后台每 15 分钟刷新实时链上快照，读侧无感知
	go rtRefresher(pctx)

	srv = &http.Server{
		Addr:    config.GlobalConfig.Listen,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		<-pctx.Done()
		log.Println("shutting down server...")
		ctx, cancel := context.WithTimeout(pctx, 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("server shutdown failed:%+v", err)
		}
		done <- struct{}{}
	}()
	return done
}

func addCors() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Accept", "Authorization", "Cache-Control", "Content-Type", "DNT", "If-Modified-Since", "Keep-Alive", "Origin", "User-Agent", "X-Requested-With"},
		},
	)
}

func checkHost(c *gin.Context) {
	if c.Request.Host != "soul-api.48.club" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
