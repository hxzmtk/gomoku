package manager

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/zqhhh/gomoku/errex"
	"github.com/zqhhh/gomoku/internal/httpserver"
)

type ClientManager struct {
	server *httpserver.Server
}

func (m *ClientManager) Init() error {
	server := httpserver.NewServer(fmt.Sprintf(":%d", httpPort))
	server.SessionCreator = func(c *httpserver.Conn) httpserver.Session { return &Session{} }
	m.server = server
	err := server.Start()
	if err != nil {
		return err
	}
	log.Infof("server in port:%d", httpPort)
	return nil
}

func (m *ClientManager) IsOnline(username string) bool {

	user, ok := manager.UserManager.users[username]
	if ok && user.Online() {
		return true
	}
	return false
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

type Session struct {
	conn *httpserver.Conn
}

func (s *Session) OnConnect(c *httpserver.Conn) {
	s.conn = c
}

func (s *Session) OnMessage(data []byte) {
	c := s.conn
	rcv, err := httpserver.Unmarshal(data)
	if err != nil {
		log.Debugf("error: %v", err)
		c.WriteMessage(&httpserver.MsgErrorAck{Msg: "不支持的协议格式"})
		return
	}
	rcvMsg, err := httpserver.DoHandle(c, rcv)
	if err != nil {
		switch e := err.(type) {
		case errex.Item:
			c.WriteMessage(&httpserver.MsgErrorAck{Msg: e.Message})
		default:
			c.WriteMessage(&httpserver.MsgErrorAck{Msg: errex.ErrFail.Message})
			log.Infof("handle error:%v", err)
		}
	} else {
		c.WriteMessage(rcvMsg)
	}
}

func (s *Session) OnClose(c *httpserver.Conn) {
	manager.UserManager.disconnect(c.Username)
	manager.RoomManager.delete(c.Username)
}
