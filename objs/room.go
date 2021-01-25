package objs

import (
	"math/rand"

	"github.com/zqhhh/gomoku/errex"
	"github.com/zqhhh/gomoku/internal/chessboard"
	"github.com/zqhhh/gomoku/internal/httpserver"
)

type Room struct {
	Id          int
	Master      *User
	Enemy       *User
	firstMove   *User
	currentMove *User
	winner      *User
	watch       map[string]*User
	chessboard  chessboard.Node
	Started     bool
}

func (m *Room) GetEnemy(user *User) *User {
	if m.Master == user {
		return m.Enemy
	}
	return m.Master
}

func (m *Room) IsFull() bool {
	if m.Master != nil && m.Enemy != nil {
		return true
	}
	return false
}

func (m *Room) reset() {
	m.chessboard.Reset()
	m.Started = false
	m.firstMove = nil
	m.currentMove = nil
	m.winner = nil
}

func (m *Room) random() {
	randId := rand.Intn(2)
	if randId == 0 {
		m.firstMove = m.Master
	} else {
		m.firstMove = m.Enemy
	}
	m.currentMove = m.firstMove
}

func (m *Room) Start() {
	m.reset()
	m.random()
	m.Started = true
	m.ntfStartGame()
}

func (m *Room) Restart() {
	m.reset()
	m.random()
	m.Started = true
	m.ntfRestartGame()
}

func (m *Room) Leave(user *User) {
	if m.Master != user && m.Enemy != user {
		return
	}
	m.reset()
	if m.Master == user { //转移房主
		m.Master = m.Enemy
	}
	m.Enemy = nil
	if m.Master != nil {
		m.Master.Ntf(&httpserver.NtfLeaveRoom{})
	}
}

func (m *Room) ntfStartGame() {
	hand := chessboard.WhiteHand
	if m.firstMove == m.Master {
		hand = chessboard.BlackHand
	}
	m.Master.Ntf(&httpserver.NtfStartGame{Hand: hand})
	m.Enemy.Ntf(&httpserver.NtfStartGame{Hand: hand.Reverse()})
}

func (m *Room) ntfRestartGame() {
	hand := chessboard.WhiteHand
	if m.firstMove == m.Master {
		hand = chessboard.BlackHand
	}
	m.Master.Ntf(&httpserver.NtfRestartGame{Hand: hand})
	m.Enemy.Ntf(&httpserver.NtfRestartGame{Hand: hand.Reverse()})
}

func (m *Room) NtfJoinRoom() {
	user := m.Enemy
	if user == nil {
		user = m.Master
	}
	m.Master.Ntf(&httpserver.NtfJoinRoom{Username: user.Username})
}

func (m *Room) ntfWalk(x, y int, hand chessboard.Hand) {
	m.Master.Ntf(&httpserver.NtfWalk{X: x, Y: y, Hand: hand})
	m.Enemy.Ntf(&httpserver.NtfWalk{X: x, Y: y, Hand: hand})
}

func (m *Room) ntfGameOver() {
	if m.Master == m.winner {
		m.Master.Ntf(&httpserver.NtfGameOver{Msg: "恭喜您获得胜利"})
		m.Enemy.Ntf(&httpserver.NtfGameOver{Msg: "您输了,请再接再厉"})
	} else {
		m.Enemy.Ntf(&httpserver.NtfGameOver{Msg: "恭喜您获得胜利"})
		m.Master.Ntf(&httpserver.NtfGameOver{Msg: "您输了,请再接再厉"})
	}
}

func (m *Room) GoSet(user *User, x, y int) error {
	if m.winner != nil {
		return errex.ErrGameOver
	}
	if m.currentMove != user {
		return errex.ErrNotCurrentYou
	}
	success := m.chessboard.Go(x, y, m.GetMyHand(user))
	if !success {
		return errex.ErrInvalidPos
	}
	m.currentMove = m.GetEnemy(user)
	m.ntfWalk(x, y, m.GetMyHand(user))
	if m.chessboard.IsWin(x, y) {
		m.ntfGameOver()
		m.winner = user
	}
	return nil
}

func (m *Room) GetMyHand(user *User) chessboard.Hand {
	if m.firstMove == user {
		return chessboard.BlackHand
	}
	return chessboard.WhiteHand
}

func NewRoom() *Room {
	return &Room{
		chessboard: chessboard.NewChessboard(15),
		watch:      make(map[string]*User),
	}
}
