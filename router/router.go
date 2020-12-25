package router

import (
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/hxzmtk/gomoku/handler/html"
	"github.com/hxzmtk/gomoku/handler/ws/human"
	"github.com/hxzmtk/gomoku/handler/ws/robot"
)

func RegisterRouter() *gin.Engine {
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	}
	engine := gin.Default()

	// load static
	engine.Static("/static", "web/static")

	// load html template
	engine.LoadHTMLFiles("web/chess.html", "web/ai.html")

	engine.GET("/", html.Index)
	engine.GET("/ai", html.AI)

	ws := engine.Group("/ws")
	ws.GET("/human", human.Handler)
	ws.GET("/ai", robot.Handler)
	return engine
}
