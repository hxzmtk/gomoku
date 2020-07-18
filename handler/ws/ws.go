package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func init() {
	Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}
