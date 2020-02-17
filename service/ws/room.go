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

func (h *Hub) CreateRoom(master *Client) (int, error) {
	h.mux.Lock()
	defer h.mux.Unlock()

	for i, _ := range h.Rooms {
		if h.Rooms[i].Master == master {
			return 0, errors.New("您已创建过房间啦")
		}
	}
	rand.Seed(time.Now().Unix())

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
	}
	return roomID, nil

}

func (h *Hub) JoinRoom(c *Client, roomID int) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	roomNumber := uint(roomID)
	if r, ok := h.Rooms[uint(roomID)]; ok {

		if r.FirstMove == nil && r.LastMove == nil {
			if h.Rooms[roomNumber].Master == nil {
				h.Rooms[roomNumber].Master = c
			} else {
				if h.Rooms[roomNumber].Master == c {
					return errors.New("您已在房间")
				}
			}

			//选择谁先手
			rand.Seed(time.Now().Unix())
			if rand.Intn(10)%2 == 0 {
				h.Rooms[roomNumber].FirstMove = c
				h.Rooms[roomNumber].LastMove = h.Rooms[roomNumber].Master
			} else {
				h.Rooms[roomNumber].FirstMove = h.Rooms[roomNumber].Master
				h.Rooms[roomNumber].LastMove = c
			}
		}
		if r.FirstMove != nil && r.LastMove != nil && r.Master != nil && r.grid == nil {
			r.grid = gomoku.InitGrid(15, 15, &gomoku.Grid{})
		}
	} else {
		return errors.New("房间不存在")
	}
	return nil
}

//落子
func (h *Hub) GoSet(c *Client, msg *RcvChessMsg) (bool, string) {
	if msg.RoomNumber <= 0 {
		return false, "无效的房间号"
	}
	roomID := uint(msg.RoomNumber)
	if r, ok := h.Rooms[roomID]; ok {
		h.Rooms[roomID].mux.Lock()
		defer h.Rooms[roomID].mux.Unlock()

		if r.isWin {
			return false, "已分出胜负了"
		}
		if r.FirstMove != c && r.LastMove != c {
			return false, "您是观战用户"
		}
		if r.NextWho != nil && r.NextWho != c {
			return false, "请等待对手落子"
		}
		if r.FirstMove == c {
			msg.IsBlack = true
			if r.grid.Set(msg.Y+1, msg.X+1, gomoku.BlackHand) {
				r.NextWho = r.LastMove
			}
			if r.grid.IsWin(msg.Y+1, msg.X+1) {
				r.isWin = true
				return true, fmt.Sprintf("黑手赢")
			}
		} else if r.LastMove == c {
			if r.grid.Set(msg.Y+1, msg.X+1, gomoku.WhiteHand) {
				r.NextWho = r.FirstMove
			}
			if r.grid.IsWin(msg.Y+1, msg.X+1) {
				r.isWin = true
				return true, fmt.Sprintf("白手赢")
			}

		}
	} else {
		return false, "房间不存在"
	}
	return true, ""
}
