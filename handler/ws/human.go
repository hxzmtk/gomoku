package ws

import (
	"github.com/bzyy/gomoku/pkg/util"
	serviceHub "github.com/bzyy/gomoku/service/hub"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Human(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &serviceHub.HumanClient{
		Conn: conn,
		Hub:  hub,
		//send: make(chan []byte, 256),
	}

	//TODO 验证生成的ID(名字)是否已存在
	clientID := util.GetRandomName()
	client.ID = clientID

	client.Hub.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump()
}
