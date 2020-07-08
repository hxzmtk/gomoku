package ws

import (
	serviceHub "github.com/bzyy/gomoku/service/hub"
	"github.com/gorilla/websocket"
	"sync"
)

var hub *serviceHub.Hub
var once sync.Once

func init() {
	once.Do(func() {
		hub = serviceHub.NewHub()
		go hub.Run()
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
