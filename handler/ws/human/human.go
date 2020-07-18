package human

import (
	"github.com/bzyy/gomoku/handler/ws"
	"github.com/bzyy/gomoku/pkg/util"
	"github.com/bzyy/gomoku/service/hub"
	"github.com/gin-gonic/gin"
	"log"
)

func Handler(c *gin.Context) {
	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		Conn: conn,
		Send: make(chan hub.IMsg, 256),
	}

	//TODO 验证生成的ID(名字)是否已存在
	clientID := util.GetRandomName()
	client.ID = clientID

	hub.Hub.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump()
}
