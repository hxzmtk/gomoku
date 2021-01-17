package httpserver

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/zqhhh/gomoku/errex"
	"github.com/zqhhh/gomoku/pkg/util"
)

type Pumper interface {
	writePump()
	readPump()
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Conn struct {
	ws       *websocket.Conn
	hub      *Hub
	Username string
	send     chan IMessage
}

func (conn Conn) GetId() int {
	return 0
}

func (conn *Conn) Start() {
	go func() {
		conn.readPump()
	}()
	go func() {
		conn.writePump()
	}()
}

func (c Conn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message.ToBytes())

			// Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	msg := <-c.send
			// 	c.ws.WriteMessage(websocket.TextMessage, msg.ToBytes())
			// }
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func (c *Conn) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Infof("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		rcv, err := Unmarshal(message)
		if err != nil {
			log.Errorf("error: %v", err)
			continue
		}
		rcvMsg, err := DoHandle(c, rcv)
		if err != nil {
			switch e := err.(type) {
			case errex.Item:
				c.WriteMessage(&MsgErrorAck{Msg: e.Message})
			default:
				c.WriteMessage(&MsgErrorAck{Msg: errex.ErrFail.Message})
				log.Infof("handle error:%v", err)
			}
		} else {
			c.WriteMessage(rcvMsg)
		}
		c.hub.broadcast <- message
	}
}

func (c *Conn) WriteMessage(msg IMessage) {
	msg.SetMsgId(msg.GetMsgId())
	c.send <- msg
}

func (c *Conn) Init() {
}

func NewConn(c *websocket.Conn, hub *Hub) *Conn {
	username := util.GetRandomName()
	conn := &Conn{ws: c,
		hub:      hub,
		Username: username,
		send:     make(chan IMessage, 1024),
	}
	hub.register <- conn
	return conn
}
