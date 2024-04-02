package manager

import (
	"math/rand"
	"sync"
	"time"

	"github.com/zqb7/gomoku/internal/session"
	"github.com/zqb7/gomoku/objs"
	"github.com/zqb7/gomoku/pkg/errex"
)

type UserManager struct {
	users map[string]*objs.User
	mux   sync.RWMutex
}

func (*UserManager) Init() error {
	return nil
}
func (m *UserManager) AddUser(user *objs.User) {
	m.users[user.Username] = user
}

func (m *UserManager) LoadUser(s *session.Session) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if err := m.reconnect(s); err == nil {
		return nil
	}
	user := objs.NewUser()
	user.Username = m.RandomName()
	s.Username = user.Username
	user.Session = s
	m.AddUser(user)
	return nil
}

func (m *UserManager) GetUser(s *session.Session) *objs.User {
	user := m.users[s.Username]
	return user
}

func (m *UserManager) reconnect(s *session.Session) error {
	username := s.Username
	user, ok := m.users[username]
	if ok {
		user.Session = s
		if session := manager.ClientManager.getWaitSession(username); session != nil {
			session.StopwaitTimer <- struct{}{}
		}
		return nil
	}
	return errex.ErrReconnect
}

func (m *UserManager) disconnect(username string) {
	delete(m.users, username)
}

func (m *UserManager) RandomName() string {
	prefixStr := "abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(prefixStr)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 3; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	str := "0123456789" + prefixStr
	bytes = []byte(str)
	for i := 0; i < 3; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func NewUserManager() *UserManager {
	return &UserManager{
		users: make(map[string]*objs.User),
	}
}
