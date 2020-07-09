package hub

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"log"
)

type IMsg interface {
	send()
	receive()
	ToBytes() []byte
}

var (
	_ IMsg = &Msg{}
)

type IContent interface {
	decode()
}

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
	RoomReset
)

type Msg struct {
	MType   msgType     `json:"m_type"`
	Content interface{} `json:"content"`
	Status  bool        `json:"status"`
	Msg     string      `json:"msg"`
	client  IClient     `json:"-"`
}

func (msg *Msg) send() {

}
func (msg *Msg) receive() {
	c := msg.client.(*HumanClient)
	switch msg.MType {
	case roomMsg:
		m := RcvRoomMsg{}
		_ = mapstructure.Decode(msg.Content, &m)

		switch m.Action {
		case RoomCreate:
			if roomNumber, err := c.Hub.CreateRoom(c); err == nil {
				m.RoomNumber = roomNumber
				msg.Status = true
			} else if roomNumber == 0 {
				msg.Msg = err.Error()
			} else {
				msg.Msg = err.Error()
			}
			msg.Content = m
			c.Send <- msg
		case RoomJoin:
			if err := c.Hub.JoinRoom(c, m.RoomNumber); err != nil {
				log.Println(err)
				msg.Msg = err.Error()
			}
			enemy := c.getEnemy().(*HumanClient)
			msg.Content = ResRoomJoinMsg{Name: enemy.ID, Action: RoomJoin}
			c.Send <- msg

			//通知对方，“我”已加入房间
			if c.Room != nil && enemy != nil {
				msg.Content = ResRoomJoinMsg{Name: c.ID, Action: RoomJoin}
				msg.Status = true
				msg.Msg = "对手加入成功"
				enemy.Send <- msg
			}
		case RoomLeave:
			if c.Room != nil {
				enemy := c.getEnemy().(*HumanClient)
				isMaster := false
				if c.Room.Master != nil && c.Room.Master == c {
					isMaster = true
				}
				if enemy != nil {
					msg.Content = ResRoomLeaveMsg{IsMaster: isMaster, Action: RoomLeave}
					enemy.Send <- msg
				}
				c.Room.LeaveRoom(c)
			}
			msg.Content = m
			msg.Status = true
			msg.Msg = "您离开房间了"
			c.Send <- msg
		case RoomStart:
			if c.Room != nil {
				var err error
				if err = c.Room.Start(c); err != nil {
					msg.Msg = err.Error()
				} else {
					msg.Status = true
				}
				m.RoomNumber = int(c.Room.ID)

				if c.Room.FirstMove == c {
					m.IsBlack = true
				}
				msg.Content = m
				msg.Msg = "SUCCESS"
				c.Send <- msg // 返回开始游戏成功

				// 通知对手， 游戏要开始了
				enemy := c.getEnemy().(*HumanClient)
				if enemy != nil && err == nil {
					if c.Room.FirstMove == enemy {
						m.IsBlack = true
					} else {
						m.IsBlack = false
					}
					msg.Content = m
					msg.Msg = "房主开始了游戏"
					enemy.Send <- msg
				}
			}
		case RoomRestart:
		case RoomReset:
			if c.Room != nil && c.Room.Master == c {
				msg.Status = true
				c.Room.GameReset()
				c.Send <- msg
			}
			msg.Status = false
			msg.Msg = "房间不存在或您不是房主"
			c.Send <- msg
		}
	case chessWalk:
		m := RcvChessMsg{}
		_ = mapstructure.Decode(msg.Content, &m)

		if c.Room.GetTarget(c) == nil {
			msg.Msg = "对手断开连接了"
			msg.Content = m
			c.Send <- msg
		}
		if c.Room == nil {
		}
		if err := c.Room.GoSet(c, &m); err == nil {
			msg.Status = true
			msg.Msg = "SUCCESS"
		} else {
			msg.Msg = err.Error()
		}
		msg.Content = m
		c.Send <- msg
		enemy := c.getEnemy().(*HumanClient)
		if enemy != nil && msg.Status {
			enemy.Send <- msg
		}
	case roomList:
		msg.Content = c.Hub.GetRooms()
		c.Send <- msg
	case clientInfoMsg:
		msg.Content = ClientInfo{Name: c.ID}
		c.Send <- msg
	}
}
func (msg *Msg) ToBytes() []byte {
	message, _ := json.Marshal(msg)
	return message
}

type MsgRoomInfo struct {
	RoomNumber uint `json:"room_number" mapstructure:"room_number"`
	IsFull     bool `json:"is_full" mapstructure:"is_full"` //是否满了
}

type MsgRoomInfoList []MsgRoomInfo

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

type ResRoomLeaveMsg struct {
	Action   RoomAction `json:"action"`
	IsMaster bool       `json:"is_master"`
}
