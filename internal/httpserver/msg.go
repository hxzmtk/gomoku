package httpserver

import (
	"encoding/json"
)

type MsgRoomInfo struct {
	Number int  `json:"number"`
	IsFull bool `json:"isFull"`
}

type msgUtil struct {
	Id int `json:"id"`
}

func (m *msgUtil) GetMsgId() int {
	return m.Id
}
func (m *msgUtil) SetMsgId(id int) {
	m.Id = id
}

func (msg msgUtil) ToBytes() []byte {
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

type MsgCreateRoomReq struct {
	msgUtil
}

type MsgCreateRoomAck struct {
	msgUtil
	Number int `json:"number"`
}
