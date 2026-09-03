package db

import (
	"sp/config"
	"sp/types"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type server struct {
	*gorm.DB
}

var Server = server{}

func init() {
	engine, err := gorm.Open(mysql.Open(config.GlobalConfig.Database.DSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := engine.DB()
	if err != nil {
		panic(err)
	}
	// 连接寿命短于 MySQL wait_timeout，避免长跑后拿到服务端已关闭的连接
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	err = engine.AutoMigrate(&types.Users{}, &types.SoulPoints{}, &types.SnapTime{})
	if err != nil {
		panic(err)
	}
	Server.DB = engine
}
