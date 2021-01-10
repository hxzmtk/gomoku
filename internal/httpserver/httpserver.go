package httpserver

type Conn interface {
	GetId() int
}


type MsgId int
const (
	MsgListRoom MsgId = iota + 1
)