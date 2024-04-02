package message

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/zqb7/gomoku/internal/chessboard"
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
	msgTypes[ntfCommonMsg] = reflect.TypeOf((*NtfCommonMsg)(nil)).Elem()
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

func GetMsgId(msg IMessage) int {
	for msgId := range msgTypes {
		if reflect.TypeOf(msg).Elem() == msgTypes[msgId] {
			return msgId
		}
	}
	return 0
}

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
