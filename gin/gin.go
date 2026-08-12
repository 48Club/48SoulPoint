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
