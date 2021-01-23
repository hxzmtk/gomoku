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
	Username string `json:"username"`
}

type MsgStartGameReq struct {
	msgUtil
	RoomId int `json:"roomId"`
}
type MsgStartGameAck struct {
	msgUtil
}
