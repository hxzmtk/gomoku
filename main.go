package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	log "github.com/sirupsen/logrus"
	"github.com/zqb7/gomoku/handler"
	"github.com/zqb7/gomoku/internal/session"
	"github.com/zqb7/gomoku/manager"
	"github.com/zqb7/network"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{DisableQuote: true})
}

func setupHttpServer(port int) {
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	}
	engine := gin.Default()

	staticFileBox := packr.NewBox("./web/static")
	engine.StaticFS("/static", http.FileSystem(staticFileBox))

	indexHtml, err := packr.NewBox("./web").FindString("index.html")
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("index").Parse(indexHtml)
	if err != nil {
		panic(err)
	}
	engine.GET("/", func(c *gin.Context) {
		tmpl.Execute(c.Writer, map[string]interface{}{
			"debug": gin.IsDebugging(),
		})
	})
	wsServer := network.NewWsServer(session.NewSession)
	engine.GET("/ws", func(c *gin.Context) { wsServer.(*network.WsServer).ServeHTTP(c.Writer, c.Request) })
	engine.Run(fmt.Sprintf("127.0.0.1:%d", port))
}

func main() {
	m := manager.Get()
	if err := m.Init(); err != nil {
		log.Infoln(err)
		return
	}
	handler.Register()

	setupHttpServer(8000)

	// 优雅关闭http服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	m.Stop()
}
