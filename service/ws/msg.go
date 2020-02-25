package ws

import (
	"errors"
	"github.com/mitchellh/mapstructure"
)

type hand uint

const (
	NilHand   hand = iota //空白
	BlackHand             //黑手
	WhiteHand             //白手
)

type msgType uint
type RoomAction uint

const (
	clientInfoMsg msgType = iota //获取连接信息
	roomMsg                      //房间消息(创建房间、加入房间)
	chessWalk                    //落子消息
	roomList                     //获取房间列表消息
)

//房间消息的动作
const (
	RoomCreate RoomAction = iota
	RoomJoin
	RoomStart
	RoomLeave
	RoomRestart
)

type WsReceive struct {
	MType   msgType     `json:"m_type"`
	Content interface{} `json:"content"`
	Status  bool        `json:"status"`
	Msg     string      `json:"msg"`
}

func (w *WsReceive) verify() error {
	if w.Content == nil {
		return nil
	}
	if _, ok := w.Content.(map[string]interface{}); !ok {
		return errors.New("必须是一个字典")
	}
	switch w.MType {
	case roomMsg:
		m := RcvRoomMsg{}
		if err := mapstructure.Decode(w.Content, &m); err != nil {
			return errors.New("格式错误")
		}
		w.Content = m

	case chessWalk:
		m := RcvChessMsg{}

		if err := mapstructure.Decode(w.Content, &m); err != nil {
			return errors.New("格式错误")
		}

		if m.RoomNumber == 0 {
			return errors.New("无效的房间编号")
		}
		w.Content = m
	case roomList:
	case clientInfoMsg:

	default:
		return errors.New("错误的消息类型")
	}
	return nil
}

type RcvRoomMsg struct {
	Action     RoomAction `json:"action" mapstructure:"action"`
	RoomNumber int        `json:"room_number" mapstructure:"room_number"` //房间编号
	IsBlack    bool       `json:"is_black" mapstructure:"is_black"`
}

type RcvChessMsg struct {
	X          int  `json:"x" mapstructure:"x"` //横坐标
	Y          int  `json:"y" mapstructure:"y"` //纵坐标
	RoomNumber int  `json:"room_number" mapstructure:"room_number"`
	IsBlack    bool `json:"is_black" mapstructure:"is_black"` //是否先手
}

type ResRoomListMsg struct {
	RoomNumber uint `json:"room_number" mapstructure:"room_number"`
	IsFull     bool `json:"is_full" mapstructure:"is_full"` //是否满了
}

type MainMsg struct {
	ID  string `json:"id"`
	Msg []byte `json:"msg"`
}

type ClientInfo struct {
	Name string `json:"name" mapstructure:"name"`
}

type ResRoomJoinMsg struct {
	Action RoomAction `json:"action"`
	Name   string     `json:"name"`
}
