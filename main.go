package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	_ "sp/config"
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
}
