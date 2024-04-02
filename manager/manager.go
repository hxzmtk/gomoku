package manager

import (
	"flag"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zqb7/gomoku/internal/message"
	"github.com/zqb7/gomoku/internal/session"
)

var httpPort int = 8000

func init() {
	flag.IntVar(&httpPort, "port", 8000, "example: -port 8000")
	flag.Parse()
}

type IModule interface {
	Init() error
}
type Manager struct {
	ClientManager *ClientManager
	UserManager   *UserManager
	RoomManager   *RoomManager
	modules       []IModule
}

func (m *Manager) init() error {
	return nil
}
func (m *Manager) appendModule(module IModule) IModule {
	m.modules = append(m.modules, module)
	return module
}

func (m *Manager) Init() error {
	if err := m.init(); err != nil {
		return err
	}
	m.ClientManager = m.appendModule(NewClientManager()).(*ClientManager)
	m.UserManager = m.appendModule(NewUserManager()).(*UserManager)
	m.RoomManager = m.appendModule(NewRoomManager()).(*RoomManager)
	for _, m := range m.modules {
		if err := m.Init(); err != nil {
			return nil
		}
	}
	return nil
}

func (m *Manager) Stop() {
}

func (m *Manager) DisConnect(s *session.Session) {
	m.ClientManager.addWaitSession(s)
	m.RoomManager.notifyEnemyMsg(s.Username, "对方掉线了")
	s.WaitTimer = time.NewTimer(10 * time.Second)
	go func() {
		select {
		case <-s.WaitTimer.C:
			m.UserManager.mux.RLock()
			if user := m.UserManager.GetUser(s); user != nil && !user.Online() {
				m.UserManager.disconnect(s.Username)
				m.RoomManager.delete(s.Username)
			}
			m.UserManager.mux.RUnlock()
		case <-s.StopwaitTimer:
			m.RoomManager.notifyEnemyMsg(s.Username, "对方已重连")
			log.Debugf("user:%s, reconnect", s.Username)
		}
		s.WaitTimer.Stop()
		m.ClientManager.delWaitSession(s)
	}()
}

func (m *Manager) Handle(s *session.Session, msg message.IMessage) (message.IMessage, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%s", debug.Stack())
		}
	}()
	handle, ok := m.ClientManager.handles[msg.GetMsgId()]
	if !ok {
		log.Errorf("handle not existed,msgId:%d", msg.GetMsgId())
		return nil, nil
	}
	msg, err := handle(s, msg)
	return msg, err
}

var manager = &Manager{modules: make([]IModule, 0)}

func init() {
	session.Manager = manager
}

func Get() *Manager {
	return manager
}

func GetRoomManager() *RoomManager {
	return manager.RoomManager
}

func GetUserManager() *UserManager {
	return manager.UserManager
}
