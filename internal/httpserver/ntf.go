package httpserver

import "github.com/zqhhh/gomoku/internal/chessboard"

const (
	_ = iota + 1000
	ntfJoinRoom
	ntfStartGame
	ntfWalk
	ntfGameOver
	ntfRestartGame
	ntfLeaveRoom
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
