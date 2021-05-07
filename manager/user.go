package manager

import (
	"sync"

	"github.com/zqhhh/gomoku/errex"
	"github.com/zqhhh/gomoku/internal/httpserver"
	"github.com/zqhhh/gomoku/objs"
)

type UserManager struct {
	users map[string]*objs.User
	mux   sync.Mutex
}

func (*UserManager) Init() error {
	return nil
}
func (m *UserManager) AddUser(user *objs.User) {
	m.users[user.Username] = user
}

func (m *UserManager) LoadUser(conn *httpserver.Conn) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if err := m.reconnect(conn); err == nil {
		return nil
	}
	user := objs.NewUser()
	user.Username = conn.Username
	user.SetConn(conn)
	m.AddUser(user)
	return nil
}

func (m *UserManager) GetUser(conn *httpserver.Conn) *objs.User {
	user := m.users[conn.Username]
	return user
}

func (m *UserManager) reconnect(conn *httpserver.Conn) error {
	username := conn.Username
	user, ok := m.users[username]
	if ok {
		user.SetConn(conn)
		if session := manager.ClientManager.getWaitSession(username); session != nil {
			session.stopwaitTimer <- struct{}{}
		}
		return nil
	}
	return errex.ErrReconnect
}

func (m *UserManager) disconnect(username string) {
	delete(m.users, username)
}

func NewUserManager() *UserManager {
	return &UserManager{
		users: make(map[string]*objs.User),
	}
}
