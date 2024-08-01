package db

import (
	"fmt"
	"sp/types"

	// _ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type server struct {
	*gorm.DB
}

var Server = server{}

func init() {
	engine, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "1234", "192.168.1.2", 3306, "sp")), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = engine.AutoMigrate(&types.Users{}, &types.SoulPoints{}, &types.SnapTime{})
	if err != nil {
		panic(err)
	}
	Server.DB = engine
}
