package gin

import (
	"context"
	"encoding/json"
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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

// rtCache 实时链上数据快照；刷新时在新对象上构建、构建完成后原子替换，
// 读侧始终拿到完整的旧快照，刷新过程不阻断查询。
type rtCache struct {
	points  []types.SoulPoints // points > 0, 按 points 降序
	byAddr  map[common.Address]types.SoulPoints
	allJSON []byte // 全量响应预序列化，请求直接写出
}

// levelMins 等级下限，升序；上限由下一级下限推导，levelOf 与 /levels 共用
var levelMins = []struct {
	level string
	min   uint64
}{
	{"Member", 1},
	{"Entry", 48},
	{"Gold", 480},
	{"Platinum", 710},
}

func levelOf(points uint64) string {
	for i := len(levelMins) - 1; i >= 0; i-- {
		if points >= levelMins[i].min {
			return levelMins[i].level
		}
	}
	return ""
}

var levelsData = func() []gin.H {
	data := make([]gin.H, 0, len(levelMins))
	for i, l := range levelMins {
		var max any // 最高级无上限
		if i+1 < len(levelMins) {
			max = levelMins[i+1].min - 1
		}
		data = append(data, gin.H{"level": l.level, "min": l.min, "max": max})
	}
	return data
}()

// handlerLevelsFunc 静态等级表，不含用户数据
func handlerLevelsFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "", "data": levelsData})
}

var rtData atomic.Pointer[rtCache]

const (
	rtRefreshInterval = 15 * time.Minute
	rtColdStartWait   = 30 * time.Second // 缓存未就绪时请求最多等待首次构建的时长
)

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

	// 添加测试钱包地址
	details = append(details, config.GlobalConfig.TestAddr...)

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

// buildAndStore 在 singleflight 内构建并立即发布，等待方超时离开也不丢结果
func buildAndStore(ctx context.Context) (any, error) {
	cache, err := buildRTCache(ctx)
	if err != nil {
		return nil, err
	}
	rtData.Store(cache)
	return cache, nil
}

// refreshRTCache singleflight 合并并发刷新，完成后原子替换快照。
func refreshRTCache(ctx context.Context) (*rtCache, error) {
	v, err, _ := sfGroup.Do("rt-refresh", func() (any, error) {
		return buildAndStore(ctx)
	})
	if err != nil {
		return nil, err
	}
	return v.(*rtCache), nil
}

// awaitRTCache 等待首次构建完成或 ctx 结束；构建不随 ctx 取消，仍在后台继续
func awaitRTCache(ctx context.Context) (*rtCache, error) {
	ch := sfGroup.DoChan("rt-refresh", func() (any, error) {
		return buildAndStore(context.Background())
	})
	select {
	case r := <-ch:
		if r.Err != nil {
			return nil, r.Err
		}
		return r.Val.(*rtCache), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func rtRefresher(ctx context.Context) {
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

func handlerIndexFunc(c *gin.Context) {
	var query types.Query
	err := c.BindQuery(&query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid query", "data": []types.SoulPoints{}})
		return
	}

	cache := rtData.Load()
	if cache == nil {
		// 首次未就绪：有界等待首次构建，并发请求由 singleflight 合并；超时返回 503 让客户端稍后重试
		wctx, cancel := context.WithTimeout(c.Request.Context(), rtColdStartWait)
		defer cancel()
		cache, err = awaitRTCache(wctx)
		if err != nil {
			if wctx.Err() != nil {
				c.Header("Retry-After", "10")
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"code": http.StatusServiceUnavailable, "message": "warming up, retry later", "data": []types.SoulPoints{}})
				return
			}
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

var (
	srv     *http.Server
	sfGroup singleflight.Group
)

func Run(pctx context.Context) chan struct{} {
	done := make(chan struct{})
	r := gin.Default()
	r.TrustedPlatform = gin.PlatformCloudflare

	r.Use(addCors(), checkHost)

	r.GET("/", handlerIndexFunc)

	r.GET("/levels", handlerLevelsFunc)

	// 后台每 15 分钟刷新实时链上快照，读侧无感知
	go rtRefresher(pctx)

	srv = &http.Server{
		Addr:              config.GlobalConfig.Listen,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second, // 需大于 rtColdStartWait
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		<-pctx.Done()
		log.Println("shutting down server...")
		// pctx 已取消，不能派生，否则 Shutdown 立即返回不等待在途请求
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
