package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadStatic(engine *gin.Engine) {
	engine.Static("/static", "web/static")
}

func LoadHtml(engine *gin.Engine) {
	engine.LoadHTMLFiles("web/chess.html")
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chess.html", gin.H{
			"debug": gin.IsDebugging(),
		})
	})
}
