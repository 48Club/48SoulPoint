package main

import (
	"context"
	"sp/gin"
	"sp/snapshot"
)

var (
	ctx = context.Background()
)

func main() {
	if err := gin.Run(ctx); err != nil {
		panic(err)

	}

}

func init() {
	go snapshot.TakeSnapshot(ctx)
}
