package objs

import (
	"time"

	"github.com/zqb7/gomoku/internal/message"
	"github.com/zqb7/gomoku/internal/session"
	"github.com/zqb7/network"
)

type User struct {
	Username   string
	Session    network.Session
	CreateTime time.Time
}

func (user *User) Ntf(msg message.IMessage) {
	user.Session.(*session.Session).WriteMessage(msg)
}

func (user *User) Online() bool {
	return !user.Session.(*session.Session).IsDisconnect
}

func NewUser() *User {
	return &User{
		CreateTime: time.Now(),
	}
}
