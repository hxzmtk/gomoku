package hub

import (
	"encoding/json"
	"errors"
	"sort"
)

const (
	MaxRoomCount = 100 //最大房间数
)

type IClient interface {
	ReadPump()
	WritePump()
}

type Hub struct {
	clients    map[string]IClient
	broadcast  chan []byte
	register   chan IClient
	unregister chan IClient
	Rooms      map[uint]*Room
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]IClient),
		register:   make(chan IClient),
		unregister: make(chan IClient),
		Rooms:      make(map[uint]*Room, MaxRoomCount),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			humanClient, ok := client.(*HumanClient)
			if !ok {
				continue
			}
			h.clients[humanClient.ID] = client
		case client := <-h.unregister:
			humanClient, ok := client.(*HumanClient)
			if !ok {
				continue
			}
			if _, ok := h.clients[humanClient.ID]; ok {
				delete(h.clients, humanClient.ID)
				close(humanClient.Send)
			}
		case message := <-h.broadcast:
			_ = message
			msg := Msg{}
			_ = json.Unmarshal(message, &msg)
			for _, client := range h.clients {
				humanClient, ok := client.(*HumanClient)
				if !ok {
					continue
				}
				select {
				case humanClient.Send <- &msg:
				default:
					close(humanClient.Send)
					delete(h.clients, humanClient.ID)
				}
			}

		}
	}
}

func (h *Hub) CreateRoom(c IClient) (roomID int, err error) {
	return 0, nil
}

func (h *Hub) RegisterClient(c IClient) {
	h.register <- c
}

func (h *Hub) UnregisterClient(c IClient) {
	h.unregister <- c
}

func (h *Hub) GetRooms() MsgRoomInfoList {
	rooms := make(MsgRoomInfoList, 0)
	for _, room := range h.Rooms {
		if room == nil {
			continue
		}
		isFull := false
		if room.Master != nil && room.Target != nil {
			isFull = true
		}
		rooms = append(rooms, MsgRoomInfo{
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

func (h *Hub) JoinRoom(c IClient, roomID int) error {
	roomNumber := uint(roomID)
	if room, ok := h.Rooms[roomNumber]; ok {
		if err := room.Join(c); err != nil {
			return err
		}
	} else {
		return errors.New("房间不存在")
	}
	return nil
}
