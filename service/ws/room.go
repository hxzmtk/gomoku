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
	FirstMove *Client //先手
	LastMove  *Client //后手
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
	if room.FirstMove != c && room.LastMove != c {
		return false, "您是观战用户"
	}
	if room.NextWho != nil && room.NextWho != c {
		return false, "请等待对手落子"
	}

	if room.FirstMove == c {
		msg.IsBlack = true
		if room.grid.Set(msg.Y+1, msg.X+1, gomoku.BlackHand) {
			room.NextWho = room.LastMove
			if room.grid.IsWin(msg.Y+1, msg.X+1) {
				room.isWin = true
				return true, fmt.Sprintf("黑手赢")
			}
		} else {
			return false, "该位置已有棋子"
		}
	} else if room.LastMove == c {
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
		room.LastMove = room.Master
	} else {
		room.FirstMove = room.Master
		room.LastMove = c
	}
	room.NextWho = room.FirstMove
}

func (room *Room) SendMessage([]byte) {

}

//返回"对手"的指针
func (room *Room) GetTarget(me *Client) *Client {
	if room.FirstMove != me {
		return room.FirstMove
	}
	if room.LastMove != me {
		return room.LastMove
	}
	return nil
}

//初始化棋盘
func (room *Room) InitGrid() {
	if room.FirstMove != nil && room.LastMove != nil && room.Master != nil && room.grid == nil {
		room.grid = gomoku.InitGrid(15, 15, &gomoku.Grid{})
	}
}

//加入房间
func (room *Room) JoinRoom(c *Client) error {

	if room.FirstMove != nil && room.LastMove != nil {
		if room.Master == nil {
			room.Master = c
		} else if room.FirstMove != c && room.LastMove != c {
			return errors.New("房间已满")
		} else if room.Master == c || (room.FirstMove == c || room.LastMove == c) {
			return errors.New("您已在房间")
		}
	}
	c.Room = room
	return nil
}

func (h *Hub) CreateRoom(master *Client) (int, error) {
	h.mux.Lock()
	defer h.mux.Unlock()

	if master.Room != nil {
		return 0, errors.New("您已创建过房间啦")
	}

	roomID := rand.Intn(1000) + 1

	if _, ok := h.Rooms[uint(roomID)]; ok {
		return 0, errors.New("房间已存在")
	} else {
		grid := gomoku.InitGrid(15, 15, &gomoku.Grid{})
		room := &Room{
			ID:     uint(roomID),
			Master: master,
			grid:   grid,
		}
		h.Rooms[uint(roomID)] = room
		master.Room = room
	}
	return roomID, nil

}

func (h *Hub) JoinRoom(c *Client, roomID int) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	roomNumber := uint(roomID)
	if room, ok := h.Rooms[roomNumber]; ok {
		if err := room.JoinRoom(c); err != nil {
			return err
		}
		room.ELectWhoFirst(c)

	} else {
		return errors.New("房间不存在")
	}
	return nil
}

func (h *Hub) GetRooms() []ResRoomListMsg {
	rooms := []ResRoomListMsg{}
	for _, room := range h.Rooms {
		isFull := false
		if room.FirstMove != nil && room.LastMove != nil {
			isFull = true
		}
		rooms = append(rooms, ResRoomListMsg{
			RoomNumber: room.ID,
			IsFull:     isFull,
		})
	}
	return rooms
}
