package manager

import (
	"github.com/zqhhh/gomoku/internal/httpserver"
	"github.com/zqhhh/gomoku/objs"
)

type UserManager struct {
	users map[string]*objs.User
}

func (UserManager) Init() error {
	return nil
}
func (m *UserManager) AddUser(user *objs.User) {
	m.users[user.Username] = user
}

func (m *UserManager) LoadUser(conn *httpserver.Conn) error {
	user := objs.NewUser()
	user.Username = conn.Username
	user.SetConn(conn)
	m.AddUser(user)
	return nil
}

func NewUserManager() *UserManager {
	return &UserManager{
		users: make(map[string]*objs.User),
	}
}
