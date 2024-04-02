package message

import "github.com/zqb7/gomoku/internal/chessboard"

const (
	_ = iota + 1000
	ntfJoinRoom
	ntfStartGame
	ntfWalk
	ntfGameOver
	ntfRestartGame
	ntfLeaveRoom
	ntfWalkWatchingUser
	ntfAskRegret
	ntfAgreeRegret
	ntfSyncWalk
	ntfCommonMsg
)

type NtfJoinRoom struct {
	msgUtil
	Username string `json:"username"`
}

type NtfStartGame struct {
	msgUtil
	Hand chessboard.Hand `json:"hand"`
}

type NtfWalk struct {
	msgUtil
	X    int             `json:"x"`
	Y    int             `json:"y"`
	Hand chessboard.Hand `json:"hand"`
}

type NtfGameOver struct {
	msgUtil
	Msg string `json:"msg"`
}

type NtfRestartGame struct {
	msgUtil
	Hand chessboard.Hand `json:"hand"`
}

type NtfLeaveRoom struct {
	msgUtil
}

type NtfWalkWatchingUser struct {
	msgUtil
	Walks  chessboard.XYS `json:"walks"`
	Latest chessboard.XY  `json:"latest"`
}

type NtfAskRegret struct {
	msgUtil
}

type NtfAgreeRegret struct {
	msgUtil
	Agree bool `json:"agree"`
}

type NtfSyncWalk struct {
	msgUtil
	Walks  chessboard.XYS `json:"walks"`
	Latest chessboard.XY  `json:"latest"`
}

type NtfCommonMsg struct {
	msgUtil
	Msg string `json:"msg"`
}
