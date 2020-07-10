package hub

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
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
	c, ok := msg.client.(*HumanClient)
	if !ok{
		return
	}
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
				msg.Msg = err.Error()
				msg.Status = false
				msg.Content = ResRoomJoinMsg{IsMaster: false, RoomNumber: m.RoomNumber, Name: "", Action: RoomJoin}
				c.Send <- msg
				return
			}
			msg.Status = true
			enemyClient := c.getEnemy()
			enemy, ok := enemyClient.(*HumanClient)
			if !ok {
				return
			}
			msg.Content = ResRoomJoinMsg{IsMaster: c.isMaster(), RoomNumber: m.RoomNumber, Name: enemy.ID, Action: RoomJoin}
			c.Send <- msg

			//通知对方，“我”已加入房间
			if c.Room != nil && enemy != nil {

				// 因为msg是一个指针，当我们修改了该指针指向的字段内容，并写入enemy.Send时，有可能c.Send在阻塞，
				// 但测试msg的内容已被更改了，会导致c.Send接收到的数据和enemy.Send收到的数据一样，所以我们copy msg
				newMsg := *msg
				newMsg.Content = ResRoomJoinMsg{IsMaster: enemy.isMaster(), RoomNumber: m.RoomNumber, Name: c.ID, Action: RoomJoin}
				newMsg.Status = true
				newMsg.Msg = "对手加入成功"
				enemy.Send <- &newMsg
			}
		case RoomLeave:
			if c.Room != nil {
				c.Room.LeaveRoom(c)
				enemyClient := c.getEnemy()
				enemy, ok := enemyClient.(*HumanClient)
				if ok {
					isMaster := false
					if c.Room.Master != nil && c.Room.Master == c {
						isMaster = true
					}
					if enemy != nil {
						newMsg := *msg
						newMsg.Content = ResRoomLeaveMsg{IsMaster: isMaster, Action: RoomLeave}
						enemy.Send <- &newMsg
					}
				}
				c.Room = nil
			}

			newMsg := *msg
			newMsg.Content = ResRoomLeaveMsg{IsMaster: false, Action: RoomLeave}
			newMsg.Status = true
			newMsg.Msg = "您离开房间了"
			//c.Send <- &newMsg
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
					newMsg := *msg
					newMsg.Content = m
					newMsg.Msg = "房主开始了游戏"
					enemy.Send <- &newMsg
				}
			}
		case RoomRestart:
			if c.Room != nil {
				if err := c.Room.Restart(c); err != nil {
					msg.Status = false
					msg.Msg = err.Error()
					c.Send <- msg
					return
				}
				msg.Status = true
				msg.Msg = "SUCCESS"
				isBlack := false
				if c.Room.FirstMove == c {
					isBlack = true
				}
				msg.Content = RcvRoomMsg{Action: RoomRestart, RoomNumber: int(c.Room.ID), IsBlack: isBlack}
				c.Send <- msg

				client := c.getEnemy()
				if client != nil {
					enemy := client.(*HumanClient)
					newMsg := *msg
					isBlack := false
					if c.Room.FirstMove == c {
						isBlack = true
					}
					newMsg.Content = RcvRoomMsg{Action: RoomRestart, RoomNumber: int(c.Room.ID), IsBlack: isBlack}
					enemy.Send <- &newMsg
				}
			}
		case RoomReset:
			if c.Room != nil && c.Room.Master == c {
				msg.Status = true
				c.Room.GameReset()
			} else {
				msg.Status = false
				msg.Msg = "房间不存在或您不是房主"
			}
			c.Send <- msg
		}
	case chessWalk:
		if c.Room == nil {
			return
		}
		m := RcvChessMsg{}
		_ = mapstructure.Decode(msg.Content, &m)

		if c.Room.GetTarget(c) == nil {
			msg.Msg = "对手断开连接了"
			msg.Content = m
			c.Send <- msg
			return
		}
		if c.Room == nil {
			return
		}
		if c.Room.FirstMove == nil {
			msg.Status = false
			if c.Room.Master == c {
				msg.Msg = "请开始游戏"
			} else {
				msg.Msg = "请等待房主开始游戏"
			}
			c.Send <- msg
			return
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
	IsMaster   bool       `json:"is_master"`
	RoomNumber int        `json:"room_number" mapstructure:"room_number"`
	Action     RoomAction `json:"action"`
	Name       string     `json:"name"`
}

type ResRoomLeaveMsg struct {
	Action   RoomAction `json:"action"`
	IsMaster bool       `json:"is_master"`
}
