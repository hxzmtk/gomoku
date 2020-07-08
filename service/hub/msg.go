package hub

import "encoding/json"

type IMsg interface {
	send()
	receive()
	ToBytes() []byte
}

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
	MType   msgType  `json:"m_type"`
	Content IContent `json:"content"`
	Status  bool     `json:"status"`
	Msg     string   `json:"msg"`
}

func (msg *Msg) send() {

}
func (msg *Msg) receive() {

}
func (msg *Msg) ToBytes() []byte {
	message, _ := json.Marshal(msg)
	return message
}

type MsgRoomInfo struct {
	RoomNumber uint `json:"room_number" mapstructure:"room_number"`
	IsFull     bool `json:"is_full" mapstructure:"is_full"` //是否满了
}

func (msg *MsgRoomInfo) decode() {
}

type MsgRoomInfoList []MsgRoomInfo

func (msg *MsgRoomInfoList) decode() {

}
