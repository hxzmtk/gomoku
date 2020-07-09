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
	Enemy      IClient         //对手
	FirstMove  IClient         //先手, 用于判断 谁是黑子 谁是白子
	chessboard chessboard.Node //棋盘
	NextWho    IClient         //下一步该谁落棋
}

// 落子
func (room *Room) GoSet(c IClient, msg *RcvChessMsg) error {
	if msg.RoomNumber <= 0 || room.ID != uint(msg.RoomNumber) {
		return errors.New("无效的房间号")
	}

	if room.isWin {
		return errors.New("已分出胜负了")
	}
	if room.Master != c && room.Enemy != c {
		return errors.New("您是观战用户")
	}
	if room.NextWho != nil && room.NextWho != c {
		return errors.New("请等待对手落子")
	}

	if room.FirstMove == c {
		msg.IsBlack = true
		if room.chessboard.Go(msg.X, msg.Y, chessboard.BlackHand) {
			room.NextWho = room.Enemy
			if room.chessboard.IsWin(msg.X, msg.Y) {
				room.isWin = true
				return errors.New("黑手赢")
			}
		} else {
			return errors.New("该位置已有棋子")
		}
	} else if room.Enemy == c {
		if room.chessboard.Go(msg.X, msg.Y, chessboard.WhiteHand) {
			room.NextWho = room.Enemy
			if room.chessboard.IsWin(msg.X, msg.Y) {
				room.isWin = true
				return errors.New("白手赢")
			}
		} else {
			return errors.New("该位置已有棋子")
		}

	}
	return nil
}

//选举谁先手
func (room *Room) electWhoFirst() {
	rand.Seed(time.Now().Unix())
	if rand.Intn(10)%2 == 0 {
		room.FirstMove = room.Master
	} else {
		room.FirstMove = room.Enemy
	}
	room.NextWho = room.FirstMove
}

// 加入房间
func (room *Room) Join(c IClient) error {
	if room.Master == c || room.Enemy == c {
		return errors.New("您已在房间")
	} else if room.Master != nil && room.Enemy != nil {
		return errors.New("房间已满")
	}
	room.Enemy = c
	switch c.(type) {
	case *HumanClient:
		client := c.(*HumanClient)
		client.Room = room
	}
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
	if room.Enemy == c {
		room.Enemy = nil
	} else {
		room.Master = nil
	}
}

//返回"对手"的指针
func (room *Room) GetTarget(me IClient) IClient {
	if room.Master != me {
		return room.Master
	}
	if room.Enemy != me {
		return room.Enemy
	}
	return nil
}

//检查房间是否为空
func (room *Room) IsEmpty() bool {
	if room.Master == nil && room.FirstMove == nil && room.Enemy == nil {
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
	if room.Enemy == nil {
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
