package httpserver

import (
	"encoding/json"

	"github.com/zqhhh/gomoku/internal/chessboard"
)

type MsgRoomInfo struct {
	RoomId int    `json:"roomId"`
	IsFull bool   `json:"isFull"`
	Master string `json:"master"`
	Enemy  string `json:"enemy"`
}

type msgUtil struct {
	MsgId int `json:"msgId"`
}

func (m *msgUtil) GetMsgId() int {
	return m.MsgId
}

func (m *msgUtil) SetMsgId(id int) {
	m.MsgId = id
}

type MsgErrorAck struct {
	msgUtil
	Msg string `json:"msg"`
}

func (msg MsgErrorAck) ToBytes() []byte {
	message, _ := json.Marshal(msg)
	return message
}

type MsgRoomListReq struct {
	msgUtil
}

type MsgRoomListAck struct {
	msgUtil
	Data []MsgRoomInfo `json:"data"`
}

type MsgCreateRoomReq struct {
	msgUtil
}

type MsgCreateRoomAck struct {
	msgUtil
	RoomId int `json:"roomId"`
}

type MsgJoinRoomReq struct {
	msgUtil
	RoomId int `json:"roomId"`
}

type MsgJoinRoomAck struct {
	msgUtil
	RoomId int `json:"roomId"`
}

type MsgChessboardWalkReq struct {
	msgUtil
	RoomId int `json:"roomId"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type MsgChessboardWalkAck struct {
	msgUtil
	X    int             `json:"x"`
	Y    int             `json:"y"`
	Hand chessboard.Hand `json:"hand"`
}

type MsgConnectReq struct {
	msgUtil
	Username string `json:"username"`
}

type MsgConnectAck struct {
	msgUtil
	Username  string          `json:"username"`
	RoomId    int             `json:"roomId"`
	MyHand    chessboard.Hand `json:"myhand"`
	Walks     chessboard.XYS  `json:"walks"`
	Latest    chessboard.XY   `json:"latest"`
	IsWatcher bool            `json:"isWatcher"`
}

type MsgStartGameReq struct {
	msgUtil
	RoomId int `json:"roomId"`
}
type MsgStartGameAck struct {
	msgUtil
}

type MsgRestartGameReq struct {
	msgUtil
	RoomId int `json:"roomId"`
}

type MsgRestartGameAck struct {
	msgUtil
}

type MsgLeaveRoomReq struct {
	msgUtil
	RoomId int `json:"roomId"`
}

type MsgLeaveRoomAck struct {
	msgUtil
}

type MsgWatchGameReq struct {
	msgUtil
	RoomId int `json:"roomId"`
}

type MsgWatchGameAck struct {
	msgUtil
	RoomId int `json:"roomId"`
}

type MsgWalkRegretReq struct {
	msgUtil
}

type MsgWalkRegretAck struct {
	msgUtil
}

type MsgAgreeRegretReq struct {
	msgUtil
	Agree bool `json:"agree"`
}

type MsgAgreeRegretAck struct {
	msgUtil
}
