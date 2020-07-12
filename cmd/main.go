package main

import (
	"context"
	"fmt"
	"github.com/bzyy/gomoku/pkg/util"
	"github.com/bzyy/gomoku/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	engine := router.RegisterRouter()

	// 优雅关闭http服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	server := &http.Server{
		Addr:    util.GetEnv("ADDR", ":8000"),
		Handler: engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	<-quit
	fmt.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}
