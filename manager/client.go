package manager

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zqhhh/gomoku/errex"
	"github.com/zqhhh/gomoku/internal/httpserver"
)

type ClientManager struct {
	server       *httpserver.Server
	waitSessions map[string]*Session
}

func (m *ClientManager) Init() error {
	server := httpserver.NewServer(fmt.Sprintf(":%d", httpPort))
	server.SessionCreator = func(c *httpserver.Conn) httpserver.Session { return NewSession() }
	m.server = server
	err := server.Start()
	if err != nil {
		return err
	}
	m.waitSessions = make(map[string]*Session)
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

func (m *ClientManager) addWaitSession(s *Session) {
	m.waitSessions[s.conn.Username] = s
}

func (m *ClientManager) delWaitSession(s *Session) {
	delete(m.waitSessions, s.conn.Username)
}

func (m *ClientManager) getWaitSession(username string) *Session {
	return m.waitSessions[username]
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

type Session struct {
	conn          *httpserver.Conn
	waitTimer     *time.Timer
	stopwaitTimer chan struct{}
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
	manager.ClientManager.addWaitSession(s)
	manager.RoomManager.notifyDisconnect(c.Username)
	s.waitTimer = time.NewTimer(10 * time.Second)
	go func() {
		select {
		case <-s.waitTimer.C:
			manager.UserManager.mux.RLock()
			if user := manager.UserManager.GetUser(c); user != nil && !user.Online() {
				manager.UserManager.disconnect(c.Username)
				manager.RoomManager.delete(c.Username)
			}
			manager.UserManager.mux.RUnlock()
		case <-s.stopwaitTimer:
			log.Debugf("user:%s, reconnect", s.conn.Username)
		}
		s.waitTimer.Stop()
		manager.ClientManager.delWaitSession(s)
	}()
}

func NewSession() *Session {
	return &Session{
		stopwaitTimer: make(chan struct{}, 1),
	}
}
