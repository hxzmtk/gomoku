package handler

import (
	"github.com/zqhhh/gomoku/internal/httpserver"
)

func HandleListRoom(conn httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	ack := &httpserver.MsgRoomListAck{}
	return ack, nil
}

func HandleCreateRoom(c httpserver.IConn, msg interface{}) (httpserver.IMessage, error) {
	ack := &httpserver.MsgCreateRoomAck{}
	return ack, nil
}

func Register() {
	httpserver.Register(httpserver.MsgListRoom, HandleListRoom)
}
