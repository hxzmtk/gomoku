package session

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zqb7/gomoku/internal/message"
	"github.com/zqb7/gomoku/pkg/errex"
	"github.com/zqb7/network"
)

var Manager IManager

type Session struct {
	conn         network.Conn
	Username     string
	IsDisconnect bool

	WaitTimer     *time.Timer
	StopwaitTimer chan struct{}
}

func (s *Session) OnConnect(conn network.Conn) {
	s.conn = conn
}

func (s *Session) OnMessage(data []byte) {
	rcv, err := message.Unmarshal(data)
	if err != nil {
		log.Debugf("error: %v", err)
		s.WriteMessage(&message.MsgErrorAck{Msg: "不支持的协议格式"})
		return
	}
	rcvMsg, err := Manager.Handle(s, rcv)
	if err != nil {
		switch e := err.(type) {
		case errex.Item:
			s.WriteMessage(&message.MsgErrorAck{Msg: e.Message})
		default:
			s.WriteMessage(&message.MsgErrorAck{Msg: errex.ErrFail.Message})
			log.Infof("handle error:%v", err)
		}
	} else {
		s.WriteMessage(rcvMsg)
	}
}

func (s *Session) OnDisConnect() {
	Manager.DisConnect(s)
}

func (s *Session) WriteMessage(msg message.IMessage) {
	if msg == nil {
		return
	}
	msg.SetMsgId(message.GetMsgId(msg))
	s.conn.Write(message.Marshal(msg))
}

func NewSession() network.Session {
	return &Session{
		StopwaitTimer: make(chan struct{}, 1),
	}
}

type IManager interface {
	Handle(s *Session, msg message.IMessage) (message.IMessage, error)
	DisConnect(s *Session)
}
