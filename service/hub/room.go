package hub

import (
	"errors"
	"github.com/bzyy/gomoku/internal/chessboard"
	"math/rand"
	"time"
)

type Room struct {
	ID         uint            //房间编号
	isWin      bool            //是否分出胜负
	Master     IClient         //房主
	Target     IClient         //对手
	FirstMove  IClient         //先手
	chessboard chessboard.Node //棋盘
	NextWho    IClient         //下一步该落棋
}

// 落子
func (room *Room) GoSet(me IClient, msg interface{}) error {
	return nil
}

//选举谁先手
func (room *Room) electWhoFirst() {
	rand.Seed(time.Now().Unix())
	if rand.Intn(10)%2 == 0 {
		room.FirstMove = room.Master
	} else {
		room.FirstMove = room.Target
	}
	room.NextWho = room.FirstMove
}

// 加入房间
func (room *Room) Join(c IClient) error {
	if room.Master == c || room.Target == c {
		return errors.New("您已在房间")
	} else if room.Master != nil && room.Target != nil {
		return errors.New("房间已满")
	}
	room.Target = c
	// c.Room = room
	return nil
}

func (room *Room) initChessboard() {
	room.chessboard = chessboard.NewChessboard(15)
}

//离开房间
func (room *Room) LeaveRoom(c IClient) {
	if room.Master == c {

		room.Master = room.GetTarget(c) //转移房主
	}
	if room.Target == c {
		room.Target = nil
	} else {
		room.Master = nil
	}
}

//返回"对手"的指针
func (room *Room) GetTarget(me IClient) IClient {
	if room.Master != me {
		return room.Master
	}
	if room.Target != me {
		return room.Target
	}
	return nil
}

//检查房间是否为空
func (room *Room) IsEmpty() bool {
	if room.Master == nil && room.FirstMove == nil && room.Target == nil {
		return true
	}
	return false
}

//开始游戏
func (room *Room) Start(c IClient) error {
	if room.Master != nil && room.Master != c {
		return errors.New("您不是房主")
	}
	if room.IsEmpty() {
		return errors.New("空的房间")
	}
	if room.Target == nil {
		return errors.New("请等待对手加入")
	}
	room.electWhoFirst()
	if !room.chessboard.IsEmpty() {
		return errors.New("游戏已经开始了")
	}
	return nil
}

//重开游戏
func (room *Room) Restart(c IClient) error {
	return room.Start(c)
}

//重置
func (room *Room) GameReset() {
	room.isWin = false
	room.chessboard.Reset()
	room.FirstMove = nil
	room.NextWho = nil
}
