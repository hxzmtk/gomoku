package human

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hxzmtk/gomoku/service/hub"
)

var (
	PongWait time.Duration = 60 * time.Second
	once     sync.Once
	_        hub.IClient = &Client{}
)

func init() {
	once.Do(func() {
		if gin.IsDebugging() {
			PongWait = 10 * time.Minute
		}
	})
}

var (
	// Time allowed to read the next pong message from the peer.
	pongWait = PongWait

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	ID       string
	Conn     *websocket.Conn
	Send     chan hub.IMsg
	Room     *hub.Room
	subject  hub.ISubject  // 订阅的主题
	observer hub.IObserver // ↑↑↑
}

func (c *Client) ReadPump() {
	defer func() {
		c.close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		//c.Hub.broadcast <- message
		msg := &Msg{
			client: c,
		}
		if err := json.Unmarshal(message, msg); err != nil {
			log.Println("非法的消息格式", err)
			continue
		}
		msg.Receive()
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message.ToBytes())

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				msg := <-c.Send
				w.Write(msg.ToBytes())
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) getEnemy() hub.IClient {
	if c.Room == nil {
		return nil
	}
	if c.Room.Master != nil && c.Room.Master != c {
		return c.Room.Master
	}
	if c.Room.Enemy != nil && c.Room.Enemy != c {
		return c.Room.Enemy
	}
	return nil
}

// 检查是否是房主
func (c *Client) isMaster() bool {
	if c.Room == nil {
		return false
	}
	if c.Room.Master != nil && c.Room.Master == c {
		return true
	}
	return false
}

// 断开连接后，自动离开房间，退订主题等
func (c *Client) close() {
	hub.Hub.UnregisterClient(c)
	if c.Room != nil {
		c.Room.LeaveRoom(c)
	}
	if c.subject != nil && c.observer != nil {
		c.subject.Detach(c.observer)
	}
	c.Conn.Close()
}

func (c *Client) isBlack() bool {
	if room := c.Room; room != nil {
		return room.FirstMove == c
	}
	return false
}

func (c *Client) GetRoom() *hub.Room {
	return c.Room
}

func (c *Client) SetRoom(room *hub.Room) {
	c.Room = room
}

func (c *Client) GetID() (clientID string) {
	return c.ID
}

func (c *Client) CloseChan() {
	close(c.Send)
}
