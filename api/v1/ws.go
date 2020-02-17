package v1

import (
	serviceWs "github.com/bzyy/gomoku/service/ws"
	"github.com/gin-gonic/gin"
)

var hub *serviceWs.Hub

func init() {
	hub = serviceWs.NewHub()
	go hub.Run()
}
func ws(c *gin.Context) {
	serviceWs.ServeWs(hub, c)
}
