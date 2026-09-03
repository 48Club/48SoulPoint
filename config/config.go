package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sp/types"
)

type Config struct {
	RPC      string                   `json:"rpc"`
	Database database                 `json:"database"`
	TestAddr []types.CalculatorDetail `json:"test_addr"` // 测试钱包地址及积分
	Listen   string                   `json:"listen"`
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
	// SP_CONFIG 可指定配置路径，便于在非工作目录（如 go test）下运行
	path := os.Getenv("SP_CONFIG")
	if path == "" {
		path = "config.json"
	}
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &GlobalConfig)
	if err != nil {
		panic(err)
	}

	// 配置中缺失/null 的积分字段补为 0，避免 nil *big.Int 调用时 panic
	for i := range GlobalConfig.TestAddr {
		d := &GlobalConfig.TestAddr[i]
		for _, p := range []**big.Int{&d.KogePoint, &d.StakePoint, &d.NftPoint, &d.BscStakePoint, &d.GovBNBPoint} {
			if *p == nil {
				*p = new(big.Int)
			}
		}
	}
}

func (db database) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.Username, db.Password, db.Host, db.Port, db.Database)
}
