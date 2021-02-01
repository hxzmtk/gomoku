package httpserver

import (
	"encoding/json"
	"errors"
	"reflect"
)

var msgTypes = make(map[int]reflect.Type)

func init() {
	msgTypes[MsgConnect] = reflect.TypeOf((*MsgConnectReq)(nil)).Elem()
	msgTypes[-MsgConnect] = reflect.TypeOf((*MsgConnectAck)(nil)).Elem()
	msgTypes[MsgError] = reflect.TypeOf((*MsgErrorAck)(nil)).Elem()
	msgTypes[MsgCreateRoom] = reflect.TypeOf((*MsgCreateRoomReq)(nil)).Elem()
	msgTypes[-MsgCreateRoom] = reflect.TypeOf((*MsgCreateRoomAck)(nil)).Elem()
	msgTypes[MsgListRoom] = reflect.TypeOf((*MsgRoomListReq)(nil)).Elem()
	msgTypes[-MsgListRoom] = reflect.TypeOf((*MsgRoomListAck)(nil)).Elem()
	msgTypes[MsgJoinRoom] = reflect.TypeOf((*MsgJoinRoomReq)(nil)).Elem()
	msgTypes[-MsgJoinRoom] = reflect.TypeOf((*MsgJoinRoomAck)(nil)).Elem()
	msgTypes[MsgChessboardWalk] = reflect.TypeOf((*MsgChessboardWalkReq)(nil)).Elem()
	msgTypes[-MsgChessboardWalk] = reflect.TypeOf((*MsgChessboardWalkAck)(nil)).Elem()
	msgTypes[MsgStartGame] = reflect.TypeOf((*MsgStartGameReq)(nil)).Elem()
	msgTypes[-MsgStartGame] = reflect.TypeOf((*MsgStartGameAck)(nil)).Elem()
	msgTypes[MsgRestartGame] = reflect.TypeOf((*MsgRestartGameReq)(nil)).Elem()
	msgTypes[-MsgRestartGame] = reflect.TypeOf((*MsgRestartGameAck)(nil)).Elem()
	msgTypes[MsgLeaveRoom] = reflect.TypeOf((*MsgLeaveRoomReq)(nil)).Elem()
	msgTypes[-MsgLeaveRoom] = reflect.TypeOf((*MsgLeaveRoomAck)(nil)).Elem()
	msgTypes[MsgWatchGame] = reflect.TypeOf((*MsgWatchGameReq)(nil)).Elem()
	msgTypes[-MsgWatchGame] = reflect.TypeOf((*MsgWatchGameAck)(nil)).Elem()
	msgTypes[MsgWalkRegret] = reflect.TypeOf((*MsgWalkRegretReq)(nil)).Elem()
	msgTypes[-MsgWalkRegret] = reflect.TypeOf((*MsgWalkRegretAck)(nil)).Elem()
	msgTypes[MsgAgreeRegret] = reflect.TypeOf((*MsgAgreeRegretReq)(nil)).Elem()
	msgTypes[-MsgAgreeRegret] = reflect.TypeOf((*MsgAgreeRegretAck)(nil)).Elem()

	// ntf
	msgTypes[ntfJoinRoom] = reflect.TypeOf((*NtfJoinRoom)(nil)).Elem()
	msgTypes[ntfStartGame] = reflect.TypeOf((*NtfStartGame)(nil)).Elem()
	msgTypes[ntfWalk] = reflect.TypeOf((*NtfWalk)(nil)).Elem()
	msgTypes[ntfGameOver] = reflect.TypeOf((*NtfGameOver)(nil)).Elem()
	msgTypes[ntfRestartGame] = reflect.TypeOf((*NtfRestartGame)(nil)).Elem()
	msgTypes[ntfLeaveRoom] = reflect.TypeOf((*NtfLeaveRoom)(nil)).Elem()
	msgTypes[ntfWalkWatchingUser] = reflect.TypeOf((*NtfWalkWatchingUser)(nil)).Elem()
	msgTypes[ntfAskRegret] = reflect.TypeOf((*NtfAskRegret)(nil)).Elem()
	msgTypes[ntfAgreeRegret] = reflect.TypeOf((*NtfAgreeRegret)(nil)).Elem()
	msgTypes[ntfSyncWalk] = reflect.TypeOf((*NtfSyncWalk)(nil)).Elem()
}

type IConn interface {
	GetId() int
}

type IMessage interface {
	GetMsgId() int
	SetMsgId(int)
}

const (
	MsgConnect = 99999
)

const (
	// 定义消息类型
	MsgError = iota
	MsgListRoom
	MsgCreateRoom
	MsgJoinRoom
	MsgChessboardWalk
	MsgStartGame
	MsgRestartGame
	MsgLeaveRoom
	MsgWatchGame
	MsgWalkRegret
	MsgAgreeRegret
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

func Marshal(msg IMessage) []byte {
	byteArr, _ := json.Marshal(msg)
	return byteArr
}

func getMsgId(msg IMessage) int {
	for msgId := range msgTypes {
		if reflect.TypeOf(msg).Elem() == msgTypes[msgId] {
			return msgId
		}
	}
	return 0
}
