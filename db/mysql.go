package db

import (
	"sp/config"
	"sp/types"

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

	err = engine.AutoMigrate(&types.Users{}, &types.SoulPoints{}, &types.SnapTime{})
	if err != nil {
		panic(err)
	}
	Server.DB = engine
}
