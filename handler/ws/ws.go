package wsHandler

import (
	serviceWs "github.com/bzyy/gomoku/service/ws"
	"sync"
)

var hub *serviceWs.Hub
var once sync.Once

func init() {
	once.Do(func() {
		hub = serviceWs.NewHub()
		go hub.Run()
	})
}
