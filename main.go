package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/zqhhh/gomoku/manager"
	"github.com/zqhhh/gomoku/model"
)

func init() {
	model.Start()
}

func main() {
	m := manager.Get()
	if err := m.Init(); err != nil {
		log.Infoln(err)
		return
	}

	// 优雅关闭http服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	m.Stop()
}
