package manager

import (
	"sync/atomic"

	"github.com/zqhhh/gomoku/errex"
	"github.com/zqhhh/gomoku/internal/chessboard"
	"github.com/zqhhh/gomoku/internal/httpserver"
	"github.com/zqhhh/gomoku/objs"
)

type RoomManager struct {
	rooms  map[int]*objs.Room
	users  map[string]int // key=username, value=roomId
	roomId int32
}

func (RoomManager) Init() error {
	return nil
}

func (m *RoomManager) newRoomId() int {
	return int(atomic.AddInt32(&m.roomId, 1))
}

func (m *RoomManager) addRoom(room *objs.Room) {
	m.rooms[room.Id] = room
}

func (m *RoomManager) ListRooms() []objs.Room {
	rooms := make([]objs.Room, 0)
	for _, room := range m.rooms {
		rooms = append(rooms, *room)
	}
	return rooms
}

func (m *RoomManager) CreateRoom(conn *httpserver.Conn) (*objs.Room, error) {
	_, ok := m.users[conn.Username]
	if ok {
		return nil, errex.ErrDupCreateRoom
	}
	newRoom := objs.NewRoom()
	newRoom.Id = m.newRoomId()
	newRoom.Master = manager.UserManager.GetUser(conn)
	m.addRoom(newRoom)
	return newRoom, nil

}

func (m *RoomManager) JoinRoom(conn *httpserver.Conn, roomId int) error {
	room, ok := m.rooms[roomId]
	if !ok {
		return errex.ErrNotExistedRoom
	}
	user := manager.UserManager.GetUser(conn)
	if room.Master == user {
		return errex.ErrInRoom
	}
	room.Enemy = manager.UserManager.GetUser(conn)
	return nil

}

func (m *RoomManager) ChessboardWalk(conn *httpserver.Conn, roomId, x, y int) (chessboard.Hand, error) {
	room, ok := m.rooms[roomId]
	if !ok {
		return chessboard.NilHand, errex.ErrNotExistedRoom
	}
	user := manager.UserManager.GetUser(conn)
	if room.Master != user && room.Enemy != user {
		return chessboard.NilHand, errex.ErrNotInRoom
	}
	enemy := room.GetEnemy(user)
	if enemy == nil {
		return 0, errex.ErrNoEnemy
	}
	if !manager.ClientManager.IsOnline(enemy.Username) {
		return chessboard.NilHand, errex.ErrNotOnline
	}
	return chessboard.NilHand, nil
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms:  make(map[int]*objs.Room),
		users:  make(map[string]int),
		roomId: 1000,
	}
}
