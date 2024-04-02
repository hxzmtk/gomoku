package objs

import (
	"math/rand"
	"time"

	"github.com/zqb7/gomoku/internal/chessboard"
	"github.com/zqb7/gomoku/internal/message"
	"github.com/zqb7/gomoku/pkg/errex"
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
	Latest      chessboard.XY
	walkRecords chessboard.XYS
	pause       bool
	pauseAt     time.Time
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

func (m *Room) IsEmpty() bool {
	return m.Master == nil && m.Enemy == nil
}

func (m *Room) reset() {
	m.chessboard.Reset()
	m.Started = false
	m.firstMove = nil
	m.currentMove = nil
	m.winner = nil
	m.walkRecords = make(chessboard.XYS, 0)
	m.Latest.X = -1
	m.Latest.Y = -1
	m.pause = false
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
		delete(m.watch, user.Username)
		return
	}
	m.reset()
	if m.Master == user { //转移房主
		m.Master = m.Enemy
	}
	m.Enemy = nil
	if m.Master != nil {
		m.Master.Ntf(&message.NtfLeaveRoom{})
	}
}

func (m *Room) JoinWatch(user *User) error {
	if m.Master == user || m.Enemy == user {
		return errex.ErrHasInRoom
	}
	if _, ok := m.watch[user.Username]; ok {
		return errex.ErrInRoom
	}
	m.watch[user.Username] = user
	user.Ntf(&message.NtfWalkWatchingUser{Walks: m.chessboard.GetState(), Latest: m.Latest})
	return nil
}

func (m *Room) CheckIsWathingUser(username string) bool {
	_, ok := m.watch[username]
	return ok
}

func (m *Room) ntfStartGame() {
	hand := chessboard.WhiteHand
	if m.firstMove == m.Master {
		hand = chessboard.BlackHand
	}
	m.Master.Ntf(&message.NtfStartGame{Hand: hand})
	m.Enemy.Ntf(&message.NtfStartGame{Hand: hand.Reverse()})
}

func (m *Room) ntfRestartGame() {
	hand := chessboard.WhiteHand
	if m.firstMove == m.Master {
		hand = chessboard.BlackHand
	}
	m.Master.Ntf(&message.NtfRestartGame{Hand: hand})
	m.Enemy.Ntf(&message.NtfRestartGame{Hand: hand.Reverse()})
}

func (m *Room) NtfJoinRoom() {
	user := m.Enemy
	if user == nil {
		user = m.Master
	}
	m.Master.Ntf(&message.NtfJoinRoom{Username: user.Username})
}

func (m *Room) ntfWalk(x, y int, hand chessboard.Hand) {
	m.Master.Ntf(&message.NtfWalk{X: x, Y: y, Hand: hand})
	m.Enemy.Ntf(&message.NtfWalk{X: x, Y: y, Hand: hand})

	walks := m.chessboard.GetState()
	for _, user := range m.watch {
		user.Ntf(&message.NtfWalkWatchingUser{Walks: walks, Latest: chessboard.XY{X: x, Y: y, Hand: hand}})
	}
}

func (m *Room) ntfGameOver() {
	if m.Master == m.winner {
		m.Master.Ntf(&message.NtfGameOver{Msg: "恭喜您获得胜利"})
		m.Enemy.Ntf(&message.NtfGameOver{Msg: "您输了,请再接再厉"})
	} else {
		m.Enemy.Ntf(&message.NtfGameOver{Msg: "恭喜您获得胜利"})
		m.Master.Ntf(&message.NtfGameOver{Msg: "您输了,请再接再厉"})
	}
}

func (m *Room) GoSet(user *User, x, y int) error {
	if m.pause {
		if time.Now().Unix()-m.pauseAt.Unix() < 10 {
			return errex.ErrPaused
		}
		m.pause = false
		user.Ntf(&message.NtfAgreeRegret{Agree: false})
	}
	if m.winner != nil {
		return errex.ErrGameOver
	}
	if m.currentMove != user {
		return errex.ErrNotCurrentYou
	}
	err := m.chessboard.Go(x, y, m.GetMyHand(user))
	if err != nil {
		return err
	}
	m.Latest.X, m.Latest.Y, m.Latest.Hand = x, y, m.GetMyHand(user)
	m.currentMove = m.GetEnemy(user)
	m.ntfWalk(x, y, m.GetMyHand(user))
	if m.chessboard.IsWin(x, y) {
		m.ntfGameOver()
		m.winner = user
	}
	m.addRecord(m.Latest)
	return nil
}

func (m *Room) GetMyHand(user *User) chessboard.Hand {
	if m.firstMove == user {
		return chessboard.BlackHand
	}
	return chessboard.WhiteHand
}

func (m *Room) GetWalkState() chessboard.XYS {
	return m.chessboard.GetState()
}

func (m *Room) addRecord(walk chessboard.XY) {
	m.walkRecords = append(m.walkRecords, walk)
}

func (m *Room) delRecords(stepCount int) {
	length := len(m.walkRecords)
	if length <= stepCount {
		m.walkRecords = make(chessboard.XYS, 0)
		return
	}
	m.walkRecords = m.walkRecords[:length-stepCount]
}

func (m *Room) AckRegret(user *User) error {
	if !(m.Master == user || m.Enemy == user) {
		return errex.ErrIsWatchingUser
	}
	if len(m.walkRecords) < 8 {
		return errex.ErrRegretWalkLess
	}
	if m.currentMove != user {
		return errex.ErrRegretWait
	}
	m.GetEnemy(user).Ntf(&message.NtfAskRegret{})
	m.pause = true
	m.pauseAt = time.Now()
	return nil
}

func (m *Room) AgreeRegret(user *User, agree bool) {
	m.pause = false

	m.GetEnemy(user).Ntf(&message.NtfAgreeRegret{Agree: agree})
	if agree {
		m.rollBack()
	}
}

func (m *Room) rollBack() {
	stepCount := 4

	length := len(m.walkRecords)
	for i := length - stepCount; i < length; i++ {
		walk := m.walkRecords[i]
		m.chessboard.Clear(walk.X, walk.Y)
	}

	m.delRecords(stepCount)

	m.Latest = m.walkRecords[len(m.walkRecords)-1]
	ntf := &message.NtfSyncWalk{
		Walks:  m.chessboard.GetState(),
		Latest: m.Latest,
	}
	m.Master.Ntf(ntf)
	m.Enemy.Ntf(ntf)
}

func NewRoom() *Room {
	return &Room{
		chessboard:  chessboard.NewChessboard(15),
		watch:       make(map[string]*User),
		walkRecords: make(chessboard.XYS, 0),
		Latest:      chessboard.XY{X: -1, Y: -1},
	}
}
