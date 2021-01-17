package handler

import (
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
	ack := &httpserver.MsgJoinRoomAck{}
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

func Register() {
	httpserver.Register(httpserver.MsgConnect, HandleConnect)
	httpserver.Register(httpserver.MsgListRoom, HandleListRoom)
	httpserver.Register(httpserver.MsgCreateRoom, HandleCreateRoom)
	httpserver.Register(httpserver.MsgJoinRoom, HandleJoinRoom)
	httpserver.Register(httpserver.MsgChessboardWalk, HandleChessboardWalk)
}
