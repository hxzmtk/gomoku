package v1

import (
	"bytes"
	"encoding/json"
	"github.com/bzyy/gomoku/service/gomoku"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func AI(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := Client{
		conn: conn,
		send: make(chan []byte),
		gird: gomoku.InitGrid(15, 15, &gomoku.Grid{}),
		ai:   gomoku.AI{},
	}

	go client.readPump()
	go client.writePump()
}

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	gird     *gomoku.Grid
	ai       gomoku.AI
	hand     gomoku.Hand
	Start    bool
	NextHand gomoku.Hand
}

const (
	writeWait      = 100 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Msg struct {
	X      int         `json:"x" mapstructure:"x"`
	Y      int         `json:"y" mapstructure:"y"`
	Hand   gomoku.Hand `json:"hand" mapstructure:"hand"`
	Action action      `json:"action" mapstructure:"action"`
}

func (c *Client) readPump() {
	defer func() {
		close(c.send)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		msg := Msg{}
		_ = json.Unmarshal(message, &msg)
		if msg.Action == actionStart {
			c.hand = gomoku.BlackHand
			c.ai.H = c.hand.Reverse()
			c.ai.Depth = 3
			msg.Hand = c.hand
			c.Start = true
			c.NextHand = c.hand
		} else if msg.Action == actionGo {
			if c.Start == false {
				msg.Action = 0
				message, _ = json.Marshal(msg)
				c.send <- message
				continue
			}
			if c.gird.Win() {
				msg.Action = actionWin
				message, _ = json.Marshal(msg)
				c.send <- message
				continue
			}
			if c.NextHand != c.hand {
				continue
			}
			msg.Hand = c.hand
			message, _ = json.Marshal(msg)
			c.send <- message

			c.gird.SetByXY(msg.X, msg.Y, c.hand)
			c.NextHand = c.ai.H
			if c.gird.Win() {
				msg.Action = actionWin
				message, _ = json.Marshal(msg)
				c.send <- message
				continue
			}
			c.ai.NegaMax(c.gird.Copy(), c.ai.H, -gomoku.INCESSANT, gomoku.INCESSANT, c.ai.Depth)
			c.gird.SetByXY(c.ai.X, c.ai.Y, c.ai.H)
			c.NextHand = c.hand
			if c.gird.Win() {
				msg.Action = actionWin
				message, _ = json.Marshal(msg)
				c.send <- message
				continue
			}

			msg.X, msg.Y = c.ai.X, c.ai.Y
			msg.Hand = c.ai.H
		}
		message, _ = json.Marshal(msg)
		c.send <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type action uint

const (
	actionNil action = iota
	actionStart
	actionGo  //落子
	actionWin //已分出胜负
)
