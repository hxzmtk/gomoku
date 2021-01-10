package manager

import (
	"github.com/zqhhh/gomoku/internal/httpserver"
)

type ClientManager struct {
}

func (ClientManager) Init() error {
	server := httpserver.NewServer(":8000")
	return server.Start()
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}
