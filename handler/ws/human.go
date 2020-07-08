package wsHandler

import (
	serviceWs "github.com/bzyy/gomoku/service/ws"
	"github.com/gin-gonic/gin"
)

func Human(c *gin.Context) {
	serviceWs.ServeWs(hub, c)
}
