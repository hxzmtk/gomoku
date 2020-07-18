package hub

import (
	"errors"
	"github.com/bzyy/gomoku/internal/chessboard"
	"math/rand"
	"time"
)

type Room struct {
	ID               uint            //房间编号
	isWin            bool            //是否分出胜负
	Master           IClient         //房主
	Enemy            IClient         //对手
	FirstMove        IClient         //先手, 用于判断 谁是黑子 谁是白子
	chessboard       chessboard.Node //棋盘
	NextWho          IClient         //下一步该谁落棋
	WatchSubject     ISubject        //观战者主题，订阅了该主题的客户端都会收到消息
	WatchSubjectChan chan IMsg       //推送订阅
	walkHistory      IWalkHistory    //下棋步骤记录，用于悔棋操作
	Pause            bool            //是否暂停
}

// 落子
func (room *Room) GoSet(c IClient, x, y int) error {

	if room.isWin {
		return errors.New("已分出胜负了")
	}
	if room.Master != c && room.Enemy != c {
		return errors.New("您是观战用户")
	}
	if room.NextWho != nil && room.NextWho != c {
		return errors.New("请等待对手落子")
	}
	if room.FirstMove == nil {
		return errors.New("请等待房主开始游戏")
	}

	if room.FirstMove == c {
		if room.chessboard.Go(x, y, chessboard.BlackHand) {
			room.nextWhoReverse()
			if room.chessboard.IsWin(x, y) {
				room.isWin = true
				//return errors.New("黑手赢")
				return nil
			}
		} else {
			return errors.New("该位置已有棋子")
		}
	} else {
		if room.chessboard.Go(x, y, chessboard.WhiteHand) {
			room.nextWhoReverse()
			if room.chessboard.IsWin(x, y) {
				room.isWin = true
				//return errors.New("白手赢")
				return nil
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
	if room.IsEmpty() {
		room.Master = c
		c.SetRoom(room)
		return nil
	}
	if room.Master == c || room.Enemy == c {
		return errors.New("您已在房间")
	} else if room.Master != nil && room.Enemy != nil {
		return errors.New("房间已满")
	}
	room.Enemy = c
	c.SetRoom(room)
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
	room.Enemy = nil
}

//返回"对手"的指针
func (room *Room) GetTarget(me IClient) IClient {
	if room.Master != nil && room.Master != me {
		return room.Master
	}
	if room.Enemy != nil && room.Enemy != me {
		return room.Enemy
	}
	return nil
}

//检查房间是否为空
func (room *Room) IsEmpty() bool {
	if room.Master == nil && room.Enemy == nil {
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
	if room.Master != nil && room.Master != c {
		return errors.New("您不是房主")
	}
	room.FirstMove = nil
	room.chessboard.Reset()
	room.isWin = false
	room.Pause = false
	room.walkHistory.Clean()
	return nil
}

//重置
func (room *Room) GameReset() {
	room.isWin = false
	room.chessboard.Reset()
	room.FirstMove = nil
	room.NextWho = nil
	room.Pause = false
	room.walkHistory.Clean()
}

//更新下一步该谁落棋
func (room *Room) nextWhoReverse() {
	nextWho := room.NextWho
	if nextWho == room.Master {
		room.NextWho = room.Enemy
	} else {
		room.NextWho = room.Master
	}
}

// 离开观战
func (room *Room) LeaveWatch(client IClient) {
}

// 悔棋
func (room *Room) Regret() error {
	if len(room.walkHistory.GetWalks()) != 3 {
		return errors.New("暂不能悔棋")
	}
	walks := room.walkHistory.GetWalks()
	for _, walk := range walks[:len(walks)-1] {
		room.chessboard.Go(walk.X, walk.Y, chessboard.NilHand)
	}
	return nil
}

func (room *Room) WhoImHand(c IClient) chessboard.Hand {
	if room.FirstMove == nil {
		return chessboard.NilHand
	}
	if room.FirstMove == c {
		return chessboard.BlackHand
	}
	return chessboard.WhiteHand
}

func (room *Room) GetWalks() (data []chessboard.XY) {
	return room.walkHistory.GetWalks()
}
func (room *Room) RecordWalk(xy chessboard.XY) {
	if room.walkHistory != nil {
		room.walkHistory.Push(xy)
	}
}

func (room *Room) GetChessBoardState() (xys chessboard.XYS) {
	if room.chessboard != nil {
		return room.chessboard.GetState()
	}
	return
}
