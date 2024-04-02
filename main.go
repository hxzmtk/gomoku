package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/zqb7/gomoku/handler"
	"github.com/zqb7/gomoku/manager"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{DisableQuote: true})
}

func main() {
	m := manager.Get()
	if err := m.Init(); err != nil {
		log.Infoln(err)
		return
	}
	handler.Register()

	// 优雅关闭http服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	m.Stop()
}
