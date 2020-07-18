package hub

import (
	"errors"
	"fmt"
	"github.com/bzyy/gomoku/internal/chessboard"
	"sort"
	"sync"
)

const (
	MaxRoomCount = 100 //最大房间数
)

type IClient interface {
	ReadPump()
	WritePump()
	GetRoom() *Room
	SetRoom(room *Room)
	GetID() string
	CloseChan()
}

var (
	Hub     *hub
	hubOnce sync.Once
)

func init() {
	hubOnce.Do(func() {
		Hub = NewHub()
		go Hub.Run()
	})
}

type hub struct {
	clients    map[string]IClient
	broadcast  chan []byte
	register   chan IClient
	unregister chan IClient
	Rooms      map[uint]*Room
}

func NewHub() *hub {
	rooms := make(map[uint]*Room)
	for i := 1; i <= MaxRoomCount; i++ {
		rooms[uint(i)] = nil
	}
	return &hub{
		clients:    make(map[string]IClient),
		register:   make(chan IClient),
		unregister: make(chan IClient),
		Rooms:      rooms,
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:

			h.clients[client.GetID()] = client
		case client := <-h.unregister:

			if _, ok := h.clients[client.GetID()]; ok {
				delete(h.clients, client.GetID())
				client.CloseChan()
			}
		case message := <-h.broadcast:
			_ = message
			for _, client := range h.clients {
				select {
				default:
					client.CloseChan()
					delete(h.clients, client.GetID())
				}
			}

		}
	}
}

func (h *hub) CreateRoom(c IClient) (roomID int, err error) {
	client := c
	room := c.GetRoom()
	if room != nil {
		return 0, errors.New("您已创建过房间啦")
	}

	for ID, room := range h.Rooms {
		if room == nil || room.IsEmpty() {
			room := &Room{
				ID:               ID,
				Master:           client,
				chessboard:       chessboard.NewChessboard(15),
				WatchSubject:     NewSubject(),
				WatchSubjectChan: make(chan IMsg, 256),
				walkHistory:      NewWalkHistory(3),
			}
			h.Rooms[ID] = room
			client.SetRoom(room)

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

func (h *hub) RegisterClient(c IClient) {
	h.register <- c
}

func (h *hub) UnregisterClient(c IClient) {
	h.unregister <- c
}

func (h *hub) GetRooms() MsgRoomInfoList {
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

func (h *hub) JoinRoom(c IClient, roomID int) error {
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

func (h *hub) GetRoomByID(roomNumber uint) *Room {
	if room, ok := h.Rooms[roomNumber]; ok {
		return room
	}
	return nil
}
