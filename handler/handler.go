package handler

import (
	"github.com/zqhhh/gomoku/internal/httpserver"
	"github.com/zqhhh/gomoku/manager"
)

var (
	m = manager.Get()
)

func handleConnect(conn httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	m.UserManager.LoadUser(conn.(*httpserver.Conn))
	return nil, nil
}

func HandleListRoom(conn httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	ack := &httpserver.MsgRoomListAck{}
	return ack, nil
}

func HandleCreateRoom(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	ack := &httpserver.MsgCreateRoomAck{}
	return ack, nil
}

func Register() {
	httpserver.Register(httpserver.MsgConnect, handleConnect)
	httpserver.Register(httpserver.MsgListRoom, HandleListRoom)
}
