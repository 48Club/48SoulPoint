package gin

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sp/db"
	"sp/types"
	"strings"
	"time"

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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid query", "data": []gin.H{}})
		return
	}
	address := common.HexToAddress(query.Address)
	tt, _ := time.Parse("20060102", time.Now().AddDate(0, 0, -48).Format("20060102"))

	log.Println(tt.Unix())
	if address == (common.Address{}) { // empty address query all addresses
		var points []types.SoulPoints

		tx := db.Server.Model(&types.SoulPoints{}).Select("user_id, users.address AS address, SUM(points) DIV 48 AS points, created").Joins("RIGHT JOIN users ON user_id = users.id").Where("created > ?", tt.Unix()).Group("user_id").Find(&points)
		if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "mysql error", "data": []gin.H{}})
			return
		}

		if tx.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "", "data": []gin.H{}}) // no points, return empty, disable cache
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "", "data": points})
		return
	}
	if !strings.EqualFold(address.Hex(), query.Address) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid address", "data": []gin.H{}})
		return
	}

	var points types.SoulPoints
	tx := db.Server.Model(&types.SoulPoints{}).Select("user_id, users.address AS address, SUM(points) DIV 48 AS points, created").Joins("RIGHT JOIN users ON user_id = users.id").Where("users.address = ? AND created > ?", address.Hex(), tt.Unix()).Find(&points)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "mysql error", "data": []gin.H{}})
		return
	}

	if tx.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "", "data": []gin.H{}}) // no points, return empty,
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "", "data": []types.SoulPoints{points}})
}

func Run(ctx context.Context) error {

	r := gin.Default()
	var store persistence.CacheStore = persistence.NewInMemoryStore(time.Hour)
	r.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Accept", "Authorization", "Cache-Control", "Content-Type", "DNT", "If-Modified-Since", "Keep-Alive", "Origin", "User-Agent", "X-Requested-With"},
		},
	))

	r.GET("/", cache.CachePage(store, time.Hour, handlerFunc))

	return r.Run(":8545")
}
