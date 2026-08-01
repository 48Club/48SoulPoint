package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	_ "sp/config"
	"sp/db"
	"sp/gin"
	"sp/snapshot"
	"syscall"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
)

func main() {
	exit := gin.Run(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	s := <-sig
	log.Printf("Signal (%v) received, stopping\n", s)

	cancel()

	<-exit

	log.Print("server exited properly")

}

func init() {
	go snapshot.TakeSnapshot(ctx)

	// 启动时预加载明细缓存；失败不阻断启动，读取时会同步兕底加载
	go func() {
		if err := db.ReloadDetailCache(); err != nil {
			log.Printf("load detail cache error: %v", err)
		}
	}()
}
