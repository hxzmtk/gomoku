package hub

type IMsg interface {
	Receive()
	ToBytes() []byte
}

type MsgRoomInfo struct {
	RoomNumber uint `json:"room_number" mapstructure:"room_number"`
	IsFull     bool `json:"is_full" mapstructure:"is_full"` //是否满了
}

type MsgRoomInfoList []MsgRoomInfo
