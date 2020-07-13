package hub

import (
	"container/list"
	"github.com/bzyy/gomoku/internal/chessboard"
)

type IWalkHistory interface {
	Push(v chessboard.XY)
	GetWalks() []chessboard.XY
	Clean()
}

type walkHistory struct {
	len  int
	list *list.List
}

func (walk *walkHistory) Push(v chessboard.XY) {
	if walk.list.Len() >= walk.len {
		walk.list.Remove(walk.list.Back())
	}
	walk.list.PushFront(v)
}

func (walk *walkHistory) GetWalks() (data []chessboard.XY) {
	e := walk.list.Front()
	for e != nil {
		value := e.Value.(chessboard.XY)
		data = append(data, chessboard.XY{
			X:    value.X,
			Y:    value.Y,
			Hand: value.Hand,
		})
		e = e.Next()
	}
	return data
}

func (walk *walkHistory) Clean() {
	walk.list.Init()
}

func NewWalkHistory(size int) IWalkHistory {
	return &walkHistory{
		len:  size,
		list: list.New(),
	}
}
