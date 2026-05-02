package gin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sp/config"
	"sp/db"
	"sp/types"
	"strings"
	"time"

	"github.com/48Club/service_agent/handler"
	"github.com/48Club/service_agent/limit"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

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
				_ = store.Set(key, resp, expire)
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

	r.Use(addCors(), handler.LimitMiddleware, checkHost)

	r.GET("/", sfCachePage(store, time.Hour, handlerFunc))

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

func init() {
	limit.Limits = limit.IPBasedRateLimiters{
		limit.NewIPBasedRateLimiter(3, time.Second*3),
		limit.NewIPBasedRateLimiter(60, time.Minute),
	}

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
