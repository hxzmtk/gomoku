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
