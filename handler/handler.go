package handler

import (
	"github.com/zqhhh/gomoku/errex"
	"github.com/zqhhh/gomoku/internal/httpserver"
	"github.com/zqhhh/gomoku/manager"
)

var (
	m = manager.Get()
)

func HandleConnect(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	conn := c.(*httpserver.Conn)
	req := msg.(*httpserver.MsgConnectReq)
	if req.Username != "" {
		conn.Username = req.Username
	}
	if err := m.UserManager.LoadUser(conn); err != nil {
		return nil, err
	}
	ack := &httpserver.MsgConnectAck{Username: conn.Username}
	room := m.RoomManager.GetRoom(conn)
	if room != nil {
		ack = &httpserver.MsgConnectAck{
			Username:  conn.Username,
			RoomId:    room.Id,
			MyHand:    room.GetMyHand(m.UserManager.GetUser(conn)),
			Walks:     room.GetWalkState(),
			Latest:    room.Latest,
			IsWatcher: room.CheckIsWathingUser(m.UserManager.GetUser(conn).Username),
		}
	}
	return ack, nil
}

func HandleListRoom(conn httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	ack := &httpserver.MsgRoomListAck{
		Data: make([]httpserver.MsgRoomInfo, 0),
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
		ack.Data = append(ack.Data, httpserver.MsgRoomInfo{
			RoomId: room.Id,
			IsFull: room.IsFull(),
			Master: masterName,
			Enemy:  enemyName,
		})
	}
	return ack, nil
}

func HandleCreateRoom(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	conn := c.(*httpserver.Conn)
	room, err := m.RoomManager.CreateRoom(conn)
	if err != nil {
		return nil, err
	}
	ack := &httpserver.MsgCreateRoomAck{RoomId: room.Id}
	return ack, nil
}

func HandleJoinRoom(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	conn := c.(*httpserver.Conn)
	req := msg.(*httpserver.MsgJoinRoomReq)
	if err := m.RoomManager.JoinRoom(conn, req.RoomId); err != nil {
		return nil, err
	}
	ack := &httpserver.MsgJoinRoomAck{
		RoomId: req.RoomId,
	}
	return ack, nil
}

func HandleChessboardWalk(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	conn := c.(*httpserver.Conn)
	req := msg.(*httpserver.MsgChessboardWalkReq)
	hand, err := m.RoomManager.ChessboardWalk(conn, req.RoomId, req.X, req.Y)
	if err != nil {
		return nil, err
	}
	ack := &httpserver.MsgChessboardWalkAck{
		X: req.X, Y: req.Y,
		Hand: hand,
	}
	return ack, nil
}

func HandleStartGame(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	conn := c.(*httpserver.Conn)
	req := msg.(*httpserver.MsgStartGameReq)
	if err := m.RoomManager.StartGame(conn, req.RoomId); err != nil {
		return nil, err
	}
	ack := &httpserver.MsgStartGameAck{}
	return ack, nil
}

func HandleRestartGame(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	req := msg.(*httpserver.MsgRestartGameReq)
	if err := m.RoomManager.RestartGame(c.(*httpserver.Conn), req.RoomId); err != nil {
		return nil, err
	}
	ack := &httpserver.MsgStartGameAck{}
	return ack, nil
}

func HandleLeaveRoom(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	req := msg.(*httpserver.MsgLeaveRoomReq)
	if err := m.RoomManager.LeaveRoom(c.(*httpserver.Conn), req.RoomId); err != nil {
		return nil, err
	}
	ack := &httpserver.MsgLeaveRoomAck{}
	return ack, nil
}

// 请求观战
func HandleWatchGame(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	req := msg.(*httpserver.MsgWatchGameReq)
	if err := m.RoomManager.WatchGame(c.(*httpserver.Conn), req.RoomId); err != nil {
		return nil, err
	}
	ack := &httpserver.MsgWatchGameAck{RoomId: req.RoomId}
	return ack, nil
}

func HandleWalkRegret(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	_ = msg.(*httpserver.MsgWalkRegretReq)
	room := m.RoomManager.GetRoom(c.(*httpserver.Conn))
	if room == nil {
		return nil, errex.ErrNotExistedRoom
	}
	user := m.UserManager.GetUser(c.(*httpserver.Conn))
	if err := room.AckRegret(user); err != nil {
		return nil, err
	}
	ack := &httpserver.MsgWalkRegretAck{}
	return ack, nil
}

func HandleAgreeRegret(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	req := msg.(*httpserver.MsgAgreeRegretReq)
	room := m.RoomManager.GetRoom(c.(*httpserver.Conn))
	if room == nil {
		return nil, errex.ErrNotExistedRoom
	}
	user := m.UserManager.GetUser(c.(*httpserver.Conn))
	if room.CheckIsWathingUser(user.Username) {
		return nil, errex.ErrIsWatchingUser
	} else if room.Master != user && room.Enemy != user {
		return nil, errex.ErrInRoom
	}
	room.AgreeRegret(user, req.Agree)
	ack := &httpserver.MsgAgreeRegretAck{}
	return ack, nil
}

func Register() {
	httpserver.Register(httpserver.MsgConnect, HandleConnect)
	httpserver.Register(httpserver.MsgListRoom, HandleListRoom)
	httpserver.Register(httpserver.MsgCreateRoom, HandleCreateRoom)
	httpserver.Register(httpserver.MsgJoinRoom, HandleJoinRoom)
	httpserver.Register(httpserver.MsgChessboardWalk, HandleChessboardWalk)
	httpserver.Register(httpserver.MsgStartGame, HandleStartGame)
	httpserver.Register(httpserver.MsgRestartGame, HandleRestartGame)
	httpserver.Register(httpserver.MsgLeaveRoom, HandleLeaveRoom)
	httpserver.Register(httpserver.MsgWatchGame, HandleWatchGame)
	httpserver.Register(httpserver.MsgWalkRegret, HandleWalkRegret)
	httpserver.Register(httpserver.MsgAgreeRegret, HandleAgreeRegret)
}
