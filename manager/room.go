package manager

import (
	"sync"
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
	mux    sync.Mutex
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

func (m *RoomManager) addUserRecord(username string, roomId int) {
	m.users[username] = roomId
}

func (m *RoomManager) deleteUserRecord(username string) {
	delete(m.users, username)
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
	newRoom := m.createRoom()
	newRoom.Master = manager.UserManager.GetUser(conn)
	m.addRoom(newRoom)
	m.addUserRecord(newRoom.Master.Username, newRoom.Id)
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

	if room.Master == nil && room.Enemy == nil {
		room.Master = user
	} else {
		room.Enemy = manager.UserManager.GetUser(conn)
	}
	room.NtfJoinRoom()
	m.addUserRecord(user.Username, room.Id)
	return nil

}

func (m *RoomManager) ChessboardWalk(conn *httpserver.Conn, roomId, x, y int) (chessboard.Hand, error) {
	room, ok := m.rooms[roomId]
	if !ok {
		return chessboard.NilHand, errex.ErrNotExistedRoom
	}
	user := manager.UserManager.GetUser(conn)
	if room.CheckIsWathingUser(user.Username) {
		return chessboard.NilHand, errex.ErrIsWatchingUser
	}
	if room.Master != user && room.Enemy != user {
		return chessboard.NilHand, errex.ErrNotInRoom
	}
	enemy := room.GetEnemy(user)
	if enemy == nil {
		return 0, errex.ErrNoEnemy
	}
	if !enemy.Online() {
		return chessboard.NilHand, errex.ErrNotOnline
	}
	return room.GetMyHand(user), room.GoSet(user, x, y)
}

func (m *RoomManager) StartGame(conn *httpserver.Conn, roomId int) error {
	room, err := m.getRoom(roomId)
	if err != nil {
		return err
	}
	if room.Started {
		return errex.ErrGameStarted
	}
	user := manager.UserManager.GetUser(conn)
	if room.Master != user {
		return errex.ErrNotRoomMaster
	}
	if room.Enemy == nil {
		return errex.ErrNoEnemy
	}
	room.Start()
	return nil
}

func (m *RoomManager) RestartGame(conn *httpserver.Conn, roomId int) error {
	room, err := m.getRoom(roomId)
	if err != nil {
		return err
	}
	user := manager.UserManager.GetUser(conn)
	if room.Master != user {
		return errex.ErrNotRoomMaster
	}
	if room.Enemy == nil {
		return errex.ErrNoEnemy
	}
	room.Restart()
	return nil
}

func (m *RoomManager) LeaveRoom(conn *httpserver.Conn, roomId int) error {
	room, err := m.getRoom(roomId)
	if err != nil {
		return err
	}
	user := manager.UserManager.GetUser(conn)
	room.Leave(user)
	m.deleteUserRecord(user.Username)
	return nil
}

func (m *RoomManager) WatchGame(conn *httpserver.Conn, roomId int) error {
	room, err := m.getRoom(roomId)
	if err != nil {
		return err
	}
	user := manager.UserManager.GetUser(conn)
	if err := room.JoinWatch(user); err != nil {
		return err
	}
	m.addUserRecord(user.Username, roomId)
	return nil
}

func (m *RoomManager) getRoom(roomId int) (*objs.Room, error) {
	room, ok := m.rooms[roomId]
	if !ok {
		return nil, errex.ErrNotExistedRoom
	}
	return room, nil
}

func (m *RoomManager) createRoom() *objs.Room {
	m.mux.Lock()
	defer m.mux.Unlock()
	for _, m := range m.rooms {
		if m.IsEmpty() {
			return m
		}
	}
	newRoom := objs.NewRoom()
	newRoom.Id = m.newRoomId()
	return newRoom
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms:  make(map[int]*objs.Room),
		users:  make(map[string]int),
		roomId: 1000,
	}
}
