package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bzyy/gomoku/internal/chessboard"
	"sort"
)

const (
	MaxRoomCount = 100 //最大房间数
)

type IClient interface {
	ReadPump()
	WritePump()
}

var (
	_ IClient = &HumanClient{}
	_ IClient = &AIClient{}
)

type Hub struct {
	clients    map[string]IClient
	broadcast  chan []byte
	register   chan IClient
	unregister chan IClient
	Rooms      map[uint]*Room
}

func NewHub() *Hub {
	rooms := make(map[uint]*Room)
	for i := 1; i <= MaxRoomCount; i++ {
		rooms[uint(i)] = nil
	}
	return &Hub{
		clients:    make(map[string]IClient),
		register:   make(chan IClient),
		unregister: make(chan IClient),
		Rooms:      rooms,
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
	client, ok := c.(*HumanClient)
	if !ok {
		return 0, errors.New("FAIL")
	}
	if client.Room != nil {
		return 0, errors.New("您已创建过房间啦")
	}

	for ID, room := range h.Rooms {
		if room == nil || room.IsEmpty() {
			room := &Room{
				ID:               ID,
				Master:           client,
				chessboard:       chessboard.NewChessboard(15),
				WatchSubject:     NewSubject(),
				WatchSubjectChan: make(chan Msg, 256),
			}
			h.Rooms[ID] = room
			client.Room = room

			// 监听chan，向观战的客户端推送消息
			go func() {
				for msg := range room.WatchSubjectChan {
					if err := room.WatchSubject.Notify(msg); err != nil {
						fmt.Println(err)
					}
				}
			}()
			return int(ID), nil
		}
	}
	return 0, errors.New("房间数已满,不能创建更多房间啦")
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
		if room.Master != nil && room.Enemy != nil {
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

func (h *Hub) GetRoomByID(roomNumber uint) *Room {
	if room, ok := h.Rooms[roomNumber]; ok {
		return room
	}
	return nil
}
