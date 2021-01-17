package objs

import "github.com/zqhhh/gomoku/internal/chessboard"

type Room struct {
	Id         int
	Master     *User
	Enemy      *User
	firstMove  *User
	watch      map[string]*User
	chessboard chessboard.Node
}

func (m *Room) GetEnemy(user *User) *User {
	if m.Master == user {
		return m.Enemy
	}
	return m.Master
}

func (m *Room)IsFull() bool {
	if m.Master != nil && m.Enemy != nil {
		return  true
	}
	return false
}

func NewRoom() *Room {
	return &Room{
		chessboard: chessboard.NewChessboard(15),
		watch:      make(map[string]*User),
	}
}
