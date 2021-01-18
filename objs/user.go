package objs

import "github.com/zqhhh/gomoku/internal/httpserver"

type User struct {
	Username string
	conn     *httpserver.Conn
}

func (user *User) SetConn(conn *httpserver.Conn) {
	user.conn = conn
}

func NewUser() *User {
	return &User{}
}

func NewUserByConn(conn *httpserver.Conn)  *User {
	return &User{
		Username: conn.Username,
		conn: conn,
	}
}
