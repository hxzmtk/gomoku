package ws

import (
	"encoding/json"
	"sync"
)

// https://github.com/gorilla/websocket/blob/master/examples/chat/hub.go

type Hub struct {
	clients   map[string]*Client
	broadcast chan MainMsg

	register   chan *Client
	unregister chan *Client

	Rooms map[uint]*Room
	mux   sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan MainMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Rooms:      make(map[uint]*Room),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
			}
		case message := <-h.broadcast:
			msg := WsReceive{}
			_ = json.Unmarshal(message.Msg, &msg)
			for _, client := range h.clients {
				select {
				case client.send <- message.Msg:
				default:
					close(client.send)
					delete(h.clients, client.ID)
				}
			}
		}
	}
}
