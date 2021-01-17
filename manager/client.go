package manager

import (
	"github.com/zqhhh/gomoku/internal/httpserver"
)

type ClientManager struct {
	server *httpserver.Server
}

func (m *ClientManager) Init() error {
	server := httpserver.NewServer(":8000")
	m.server = server
	return server.Start()
}

func (m *ClientManager) IsOnline(username string) bool {
	return m.server.CheckOnline(username)
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}
