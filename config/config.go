package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	RPC      string   `json:"rpc"`
	Database database `json:"database"`
	Listen   string   `json:"listen"`
}

type database struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

var GlobalConfig = Config{}

func init() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &GlobalConfig)
	if err != nil {
		panic(err)
	}
	log.Println(GlobalConfig)
}

func (db database) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.Username, db.Password, db.Host, db.Port, db.Database)

}
