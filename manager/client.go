package manager

import (
	"github.com/zqb7/gomoku/internal/message"
	"github.com/zqb7/gomoku/internal/session"
	"github.com/zqb7/network"
)

type HandleFunc func(c network.Session, msg interface{}) (message.IMessage, error)
type ClientManager struct {
	waitSessions map[string]*session.Session

	handles map[int]HandleFunc
}

func (m *ClientManager) Init() error {
	m.waitSessions = make(map[string]*session.Session)
	return nil
}

func (m *ClientManager) IsOnline(username string) bool {

	user, ok := manager.UserManager.users[username]
	if ok && user.Online() {
		return true
	}
	return false
}

func (m *ClientManager) addWaitSession(s *session.Session) {
	m.waitSessions[s.Username] = s
}

func (m *ClientManager) delWaitSession(s *session.Session) {
	delete(m.waitSessions, s.Username)
}

func (m *ClientManager) getWaitSession(username string) *session.Session {
	return m.waitSessions[username]
}

func (m *ClientManager) DisConnect() {

}

func (m *ClientManager) RegisterHandle(msgID int, f HandleFunc) {
	m.handles[msgID] = f
}

func NewClientManager() *ClientManager {
	return &ClientManager{handles: make(map[int]HandleFunc)}
}
