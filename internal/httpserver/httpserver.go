package httpserver

import (
	"encoding/json"
	"errors"
	"reflect"
)

var msgTypes = make(map[int]reflect.Type)

func init() {
	msgTypes[MsgCreateRoom] = reflect.TypeOf((*MsgCreateRoomReq)(nil)).Elem()
	msgTypes[MsgListRoom] = reflect.TypeOf((*MsgRoomListReq)(nil)).Elem()
}

type IConn interface {
	GetId() int
}

type IMessage interface {
	GetMsgId() int
	ToBytes() []byte
	SetMsgId(int)
}

const (
	MsgConnect =  99999
)

const (
	// 定义消息类型
	MsgListRoom = iota + 1
	MsgCreateRoom
)

var (
	_ IConn = Conn{}
)

var (
	_ IMessage = &MsgRoomListAck{}
	_ IMessage = &MsgCreateRoomAck{}
)

type MessageFrame struct {
	MsgId int `json:"msgId"`
	Body  json.RawMessage
}

func Unmarshal(data []byte) (IMessage, error) {
	frame := &MessageFrame{}
	err := json.Unmarshal(data, frame)
	if err != nil {
		return nil, err
	}
	if msgType, ok := msgTypes[frame.MsgId]; ok {
		body := reflect.New(msgType).Interface().(IMessage)
		err := json.Unmarshal(frame.Body, body)
		body.SetMsgId(frame.MsgId)
		return body, err
	}
	return nil, errors.New("msgId error")
}
