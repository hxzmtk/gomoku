package v1

import (
	"github.com/gin-gonic/gin"
)

func LoadV1(engine *gin.Engine) {
	v1 := engine.Group("/v1")
	v1.GET("/ws", ws)
}
