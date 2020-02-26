package ws

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bzyy/gomoku/service/gomoku"
)

type Room struct {
	ID        uint    //房间编号
	isWin     bool    //是否分出胜负
	Master    *Client //房主
	Target    *Client //对手
	FirstMove *Client //谁先手
	grid      *gomoku.Grid
	mux       sync.RWMutex
	NextWho   *Client //下一步该谁落子
}

//落子
func (room *Room) GoSet(c *Client, msg *RcvChessMsg) (bool, string) {
	if msg.RoomNumber <= 0 || room.ID != uint(msg.RoomNumber) {
		return false, "无效的房间号"
	}
	room.mux.Lock()
	defer room.mux.Unlock()
	if room.isWin {
		return false, "已分出胜负了"
	}
	if room.FirstMove != c && room.Target != c {
		return false, "您是观战用户"
	}
	if room.NextWho != nil && room.NextWho != c {
		return false, "请等待对手落子"
	}

	if room.FirstMove == c {
		msg.IsBlack = true
		if room.grid.Set(msg.Y+1, msg.X+1, gomoku.BlackHand) {
			room.NextWho = room.Target
			if room.grid.IsWin(msg.Y+1, msg.X+1) {
				room.isWin = true
				return true, fmt.Sprintf("黑手赢")
			}
		} else {
			return false, "该位置已有棋子"
		}
	} else if room.Target == c {
		if room.grid.Set(msg.Y+1, msg.X+1, gomoku.WhiteHand) {
			room.NextWho = room.FirstMove
			if room.grid.IsWin(msg.Y+1, msg.X+1) {
				room.isWin = true
				return true, fmt.Sprintf("白手赢")
			}
		} else {
			return false, "该位置已有棋子"
		}

	}
	return true, ""
}

//随机选举谁先手
func (room *Room) ELectWhoFirst(c *Client) {
	rand.Seed(time.Now().Unix())
	if rand.Intn(10)%2 == 0 {
		room.FirstMove = c
	} else {
		room.FirstMove = room.Master
	}
	room.NextWho = room.FirstMove
}

func (room *Room) SendMessage([]byte) {

}

//返回"对手"的指针
func (room *Room) GetTarget(me *Client) *Client {
	if room.Master != me {
		return room.Master
	}
	if room.Target != me {
		return room.Target
	}
	return nil
}

//初始化棋盘
func (room *Room) InitGrid() {
	if room.FirstMove != nil && room.Target != nil && room.Master != nil && room.grid == nil {
		room.grid = gomoku.InitGrid(15, 15, &gomoku.Grid{})
	}
}

//加入房间
func (room *Room) JoinRoom(c *Client) error {
	if room.Target == c || room.Master == c {
		return errors.New("您已在房间")
	} else if room.Master != c && room.Target != nil {
		return errors.New("房间已满")
	}
	room.Target = c
	c.Room = room
	return nil
}

//离开房间
func (room *Room) LeaveRoom(c *Client) {
	if room.Master == c {
		room.Master = room.GetTarget(c) //转移房主
	}
	if room.Target == c {
		room.Target = nil
	} else {
		room.Master = nil
	}
	c.Room = nil
}

//检查房间是否为空
func (room *Room) IsEmpty() bool {
	if room.Master == nil && room.FirstMove == nil && room.Target == nil {
		return true
	}
	return false
}

//重置
func (room *Room) GameReset() {
	room.isWin = false
	room.grid.Reset()
	room.FirstMove = nil
	room.NextWho = nil
}

//开始游戏
func (room *Room) Start(c *Client) error {
	if room.Master != nil && room.Master != c {
		return errors.New("您不是房主")
	}
	if room.IsEmpty() {
		return errors.New("空的房间")
	}
	if room.Target == nil {
		return errors.New("请等待对手加入")
	}
	room.ELectWhoFirst(c)
	if !room.grid.IsEmpty() {
		return errors.New("游戏已经开始了")
	}
	return nil
}

//重开游戏
func (room *Room) Restart(c *Client) error {
	return room.Start(c)
}
