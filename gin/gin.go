package gin

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sp/config"
	"sp/db"
	"sp/types"
	"strings"
	"time"

	"github.com/48Club/service_agent/cloudflare"
	"github.com/48Club/service_agent/handler"
	"github.com/48Club/service_agent/limit"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
		points  []types.SoulPoints                     // mysql scan result
		res     interface{}                            // response
		dbQuery = db.Server.Model(&types.SoulPoints{}) // mysql query
		errCode = http.StatusInternalServerError       // response code
		address = common.HexToAddress(query.Address)
		tt, _   = time.Parse("20060102", time.Now().AddDate(0, 0, -48).Format("20060102"))
	)

	if query.Address == "" {
		dbQuery = dbQuery.Select("user_id, users.address AS address, SUM(points) DIV 48 AS points, created").Joins("RIGHT JOIN users ON user_id = users.id").Where("created > ?", tt.Unix()).Group("user_id").Order("SUM(points) DIV 48 DESC")
		errCode = http.StatusOK
	} else {
		if !strings.EqualFold(address.Hex(), query.Address) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid address", "data": []types.SoulPoints{}})
			return
		}
		if query.Detail {
			dbQuery = dbQuery.Select("user_id, users.address AS address, points, created")
		} else {
			dbQuery = dbQuery.Select("user_id, users.address AS address, SUM(points) DIV 48 AS points, created").Group("user_id")
		}
		dbQuery = dbQuery.Joins("RIGHT JOIN users ON user_id = users.id").Where("users.address = ? AND created > ?", address.Hex(), tt.Unix())
	}

	tx := dbQuery.Find(&points)

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
				BscStake: point.BscStakePoint,
			})
		}
		details.Points /= 48
		res = details
	} else {
		res = points
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "", "data": res})
}

var srv *http.Server

func Run(pctx context.Context) chan struct{} {
	done := make(chan struct{})
	r := gin.Default()
	cloudflare.SetRemoteAddr(r)
	store := persistence.NewInMemoryStore(time.Hour)

	r.Use(addCors(), handler.LimitMiddleware, checkHost)

	r.GET("/", cache.CachePageWithoutHeader(store, time.Hour, handlerFunc))

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
		if err := limit.Limits.SaveCache(); err != nil {
			log.Fatalf("save to cache failed:%+v", err)
		}
		done <- struct{}{}
	}()
	return done
}

func init() {
	limit.Limits = limit.IPBasedRateLimiters{
		limit.NewIPBasedRateLimiter(3, time.Second*3, "3s"),
		limit.NewIPBasedRateLimiter(60, time.Minute, "1m"),
		limit.NewIPBasedRateLimiter(3600, time.Hour*1, "1h"),
	}
	if err := limit.Limits.LoadFromCache(); err != nil {
		log.Fatalf("load from cache failed:%+v", err)
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
