package ws

import (
	"github.com/bzyy/gomoku/pkg/util"
	serviceHub "github.com/bzyy/gomoku/service/hub"
	"github.com/gin-gonic/gin"
	"log"
)

func Human(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &serviceHub.HumanClient{
		Conn: conn,
		Send: make(chan serviceHub.IMsg, 256),
	}

	//TODO 验证生成的ID(名字)是否已存在
	clientID := util.GetRandomName()
	client.ID = clientID

	serviceHub.Hub.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump()
}
