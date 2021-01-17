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

func (msg *msgUtil) ToBytes() []byte {
	message, _ := json.Marshal(msg)
	return message
}

type MsgErrorAck struct {
	msgUtil
	Msg string `json:"msg"`
}

func (MsgErrorAck) GetMsgId() int {
	return -MsgError
}

func (msg MsgErrorAck) ToBytes() []byte {
	message, _ := json.Marshal(msg)
	return message
}

type MsgRoomListReq struct {
	msgUtil
}

func (MsgRoomListReq) GetMsgId() int {
	return MsgListRoom
}

type MsgRoomListAck struct {
	msgUtil
	Data []MsgRoomInfo `json:"data"`
}

func (MsgRoomListAck) GetMsgId() int {
	return -MsgListRoom
}

func (msg MsgRoomListAck) ToBytes() []byte {
	message, _ := json.Marshal(msg)
	return message
}

type MsgCreateRoomReq struct {
	msgUtil
}

func (MsgCreateRoomReq) GetMsgId() int {
	return MsgCreateRoom
}

type MsgCreateRoomAck struct {
	msgUtil
	RoomId int `json:"roomId"`
}

func (MsgCreateRoomAck) GetMsgId() int {
	return -MsgCreateRoom
}
func (msg MsgCreateRoomAck) ToBytes() []byte {
	message, _ := json.Marshal(msg)
	return message
}

type MsgJoinRoomReq struct {
	msgUtil
	RoomId int `json:"roomId"`
}

type MsgJoinRoomAck struct {
	msgUtil
}

func (MsgJoinRoomAck) GetMsgId() int {
	return -MsgJoinRoom
}

func (msg MsgJoinRoomAck) ToBytes() []byte {
	message, _ := json.Marshal(msg)
	return message
}

type MsgChessboardWalkReq struct {
	msgUtil
	RoomId int `json:"roomId"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

func (MsgChessboardWalkReq) GetMsgId() int {
	return MsgChessboardWalk
}

type MsgChessboardWalkAck struct {
	msgUtil
	X    int             `json:"x"`
	Y    int             `json:"y"`
	Hand chessboard.Hand `json:"hand"`
}

func (MsgChessboardWalkAck) GetMsgId() int {
	return -MsgChessboardWalk
}

func (msg MsgChessboardWalkAck) ToBytes() []byte {
	message, _ := json.Marshal(msg)
	return message
}

type MsgConnectReq struct {
	msgUtil
	Username string `json:"username"`
}

func (MsgConnectReq) GetMsgId() int {
	return MsgConnect
}

type MsgConnectAck struct {
	msgUtil
	Username string `json:"username"`
}

func (MsgConnectAck) GetMsgId() int {
	return -MsgConnect
}

func (msg MsgConnectAck) ToBytes() []byte {
	message, _ := json.Marshal(msg)
	return message
}
