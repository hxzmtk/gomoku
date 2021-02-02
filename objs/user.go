package objs

import (
	"time"

	"github.com/zqhhh/gomoku/internal/httpserver"
)

type User struct {
	Username   string
	conn       *httpserver.Conn
	CreateTime time.Time
}

func (user *User) SetConn(conn *httpserver.Conn) {
	user.conn = conn
}

func (user *User) Ntf(msg httpserver.IMessage) {
	user.conn.WriteMessage(msg)
}

func (user *User) Online() bool {
	return user.conn.Online()
}

func NewUser() *User {
	return &User{
		CreateTime: time.Now(),
	}
}

func NewUserByConn(user *User, conn *httpserver.Conn) *User {
	return &User{
		Username:   conn.Username,
		conn:       conn,
		CreateTime: user.CreateTime,
	}
}
