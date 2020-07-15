package html

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "chess.html", gin.H{
		"debug": gin.IsDebugging(),
	})
}
func AI(c *gin.Context) {
	c.HTML(http.StatusOK, "ai.html", gin.H{
		"debug": gin.IsDebugging(),
	})
}
