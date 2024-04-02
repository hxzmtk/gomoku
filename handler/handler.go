package handler

import (
	"github.com/zqb7/gomoku/internal/message"
	"github.com/zqb7/gomoku/internal/session"
	"github.com/zqb7/gomoku/manager"
	"github.com/zqb7/gomoku/pkg/errex"
	"github.com/zqb7/network"
)

var (
	m = manager.Get()
)

func HandleConnect(c network.Session, msg interface{}) (message.IMessage, error) {
	s := c.(*session.Session)
	req := msg.(*message.MsgConnectReq)
	if req.Username != "" {
		s.Username = req.Username
	}
	if err := m.UserManager.LoadUser(s); err != nil {
		return nil, err
	}
	ack := &message.MsgConnectAck{Username: s.Username}
	room := m.RoomManager.GetRoom(s)
	if room != nil {
		ack = &message.MsgConnectAck{
			Username:  s.Username,
			RoomId:    room.Id,
			MyHand:    room.GetMyHand(m.UserManager.GetUser(s)),
			Walks:     room.GetWalkState(),
			Latest:    room.Latest,
			IsWatcher: room.CheckIsWathingUser(m.UserManager.GetUser(s).Username),
		}
	}
	return ack, nil
}

func HandleListRoom(conn network.Session, msg interface{}) (message.IMessage, error) {
	ack := &message.MsgRoomListAck{
		Data: make([]message.MsgRoomInfo, 0),
	}
	rooms := m.RoomManager.ListRooms()
	for _, room := range rooms {
		masterName, enemyName := "", ""
		if room.Master != nil {
			masterName = room.Master.Username
		}
		if room.Enemy != nil {
			enemyName = room.Enemy.Username
		}
		ack.Data = append(ack.Data, message.MsgRoomInfo{
			RoomId: room.Id,
			IsFull: room.IsFull(),
			Master: masterName,
			Enemy:  enemyName,
		})
	}
	return ack, nil
}

func HandleCreateRoom(c network.Session, msg interface{}) (message.IMessage, error) {
	s := c.(*session.Session)
	room, err := m.RoomManager.CreateRoom(s)
	if err != nil {
		return nil, err
	}
	ack := &message.MsgCreateRoomAck{RoomId: room.Id}
	return ack, nil
}

func HandleJoinRoom(c network.Session, msg interface{}) (message.IMessage, error) {
	s := c.(*session.Session)
	req := msg.(*message.MsgJoinRoomReq)
	if err := m.RoomManager.JoinRoom(s, req.RoomId); err != nil {
		return nil, err
	}
	ack := &message.MsgJoinRoomAck{
		RoomId: req.RoomId,
	}
	return ack, nil
}

func HandleChessboardWalk(c network.Session, msg interface{}) (message.IMessage, error) {
	s := c.(*session.Session)
	req := msg.(*message.MsgChessboardWalkReq)
	hand, err := m.RoomManager.ChessboardWalk(s, req.RoomId, req.X, req.Y)
	if err != nil {
		return nil, err
	}
	ack := &message.MsgChessboardWalkAck{
		X: req.X, Y: req.Y,
		Hand: hand,
	}
	return ack, nil
}

func HandleStartGame(c network.Session, msg interface{}) (message.IMessage, error) {
	s := c.(*session.Session)
	req := msg.(*message.MsgStartGameReq)
	if err := m.RoomManager.StartGame(s, req.RoomId); err != nil {
		return nil, err
	}
	ack := &message.MsgStartGameAck{}
	return ack, nil
}

func HandleRestartGame(c network.Session, msg interface{}) (message.IMessage, error) {
	s := c.(*session.Session)
	req := msg.(*message.MsgRestartGameReq)
	if err := m.RoomManager.RestartGame(s, req.RoomId); err != nil {
		return nil, err
	}
	ack := &message.MsgStartGameAck{}
	return ack, nil
}

func HandleLeaveRoom(c network.Session, msg interface{}) (message.IMessage, error) {
	s := c.(*session.Session)
	req := msg.(*message.MsgLeaveRoomReq)
	if err := m.RoomManager.LeaveRoom(s, req.RoomId); err != nil {
		return nil, err
	}
	ack := &message.MsgLeaveRoomAck{}
	return ack, nil
}

// 请求观战
func HandleWatchGame(c network.Session, msg interface{}) (message.IMessage, error) {
	s := c.(*session.Session)
	req := msg.(*message.MsgWatchGameReq)
	if err := m.RoomManager.WatchGame(s, req.RoomId); err != nil {
		return nil, err
	}
	ack := &message.MsgWatchGameAck{RoomId: req.RoomId}
	return ack, nil
}

func HandleWalkRegret(c network.Session, msg interface{}) (message.IMessage, error) {
	s := c.(*session.Session)
	_ = msg.(*message.MsgWalkRegretReq)
	room := m.RoomManager.GetRoom(s)
	if room == nil {
		return nil, errex.ErrNotExistedRoom
	}
	user := m.UserManager.GetUser(s)
	if err := room.AckRegret(user); err != nil {
		return nil, err
	}
	ack := &message.MsgWalkRegretAck{}
	return ack, nil
}

func HandleAgreeRegret(c network.Session, msg interface{}) (message.IMessage, error) {
	s := c.(*session.Session)
	req := msg.(*message.MsgAgreeRegretReq)
	room := m.RoomManager.GetRoom(s)
	if room == nil {
		return nil, errex.ErrNotExistedRoom
	}
	user := m.UserManager.GetUser(s)
	if room.CheckIsWathingUser(user.Username) {
		return nil, errex.ErrIsWatchingUser
	} else if room.Master != user && room.Enemy != user {
		return nil, errex.ErrInRoom
	}
	room.AgreeRegret(user, req.Agree)
	ack := &message.MsgAgreeRegretAck{}
	return ack, nil
}

func Register() {
	m.ClientManager.RegisterHandle(message.MsgConnect, HandleConnect)
	m.ClientManager.RegisterHandle(message.MsgListRoom, HandleListRoom)
	m.ClientManager.RegisterHandle(message.MsgCreateRoom, HandleCreateRoom)
	m.ClientManager.RegisterHandle(message.MsgJoinRoom, HandleJoinRoom)
	m.ClientManager.RegisterHandle(message.MsgChessboardWalk, HandleChessboardWalk)
	m.ClientManager.RegisterHandle(message.MsgStartGame, HandleStartGame)
	m.ClientManager.RegisterHandle(message.MsgRestartGame, HandleRestartGame)
	m.ClientManager.RegisterHandle(message.MsgLeaveRoom, HandleLeaveRoom)
	m.ClientManager.RegisterHandle(message.MsgWatchGame, HandleWatchGame)
	m.ClientManager.RegisterHandle(message.MsgWalkRegret, HandleWalkRegret)
	m.ClientManager.RegisterHandle(message.MsgAgreeRegret, HandleAgreeRegret)
}
