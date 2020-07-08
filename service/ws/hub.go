package wsHandler

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"

	"github.com/bzyy/gomoku/service/gomoku"
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

const MaxRoomCount = 100 //最大房间数

func NewHub() *Hub {
	rooms := make(map[uint]*Room)

	//生成一定数量的房间
	for i := 1; i <= MaxRoomCount; i++ {
		rooms[uint(i)] = nil
	}

	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan MainMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Rooms:      rooms,
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

func (h *Hub) CreateRoom(master *Client) (int, error) {
	h.mux.Lock()
	defer h.mux.Unlock()

	if master.Room != nil {
		return 0, errors.New("您已创建过房间啦")
	}

	for ID, room := range h.Rooms {
		if room == nil || room.IsEmpty() {
			grid := gomoku.InitGrid(15, 15, &gomoku.Grid{})
			room := &Room{
				ID:     ID,
				Master: master,
				grid:   grid,
			}
			h.Rooms[ID] = room
			master.Room = room
			return int(ID), nil
		}
	}
	return 0, errors.New("房间数已满,不能创建更多房间啦")
}

func (h *Hub) JoinRoom(c *Client, roomID int) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	roomNumber := uint(roomID)
	if room, ok := h.Rooms[roomNumber]; ok {
		if err := room.JoinRoom(c); err != nil {
			return err
		}
	} else {
		return errors.New("房间不存在")
	}
	return nil
}

func (h *Hub) GetRooms() []ResRoomListMsg {
	rooms := []ResRoomListMsg{}
	for _, room := range h.Rooms {
		if room == nil {
			continue
		}
		isFull := false
		if room.Master != nil && room.Target != nil {
			isFull = true
		}
		rooms = append(rooms, ResRoomListMsg{
			RoomNumber: room.ID,
			IsFull:     isFull,
		})
	}

	//房间号升序排列
	sort.Slice(rooms, func(i, j int) bool {
		if rooms[i].RoomNumber > rooms[j].RoomNumber {
			return false
		}
		return true
	})
	return rooms
}
