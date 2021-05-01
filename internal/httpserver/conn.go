package httpserver

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Pumper interface {
	writePump()
	readPump()
}

type Session interface {
	OnConnect(*Conn)
	OnMessage([]byte)
	OnClose(*Conn)
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Second

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
	Username string
	send     chan IMessage
	closed   bool
	Session  Session
}

func (conn Conn) Online() bool {
	return !conn.closed
}

func (conn Conn) GetId() int {
	return 0
}

func (c *Conn) Start() {
	go func() {
		c.readPump()
	}()
	go func() {
		c.writePump()
	}()
	c.Session.OnConnect(c)
}

func (c *Conn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
		c.Session.OnClose(c)
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
			w.Write(Marshal(message))

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
		c.ws.Close()
		c.closed = true
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
		c.Session.OnMessage(message)
	}
}

func (c *Conn) WriteMessage(msg IMessage) {
	if msg == nil {
		return
	}
	msg.SetMsgId(getMsgId(msg))
	c.send <- msg
}

func (c *Conn) Init() {
}

func NewConn(c *websocket.Conn, username string, sessionCreator func(*Conn) Session) *Conn {
	conn := &Conn{ws: c,
		Username: username,
		send:     make(chan IMessage, 1024),
	}
	conn.Session = sessionCreator(conn)
	return conn
}
