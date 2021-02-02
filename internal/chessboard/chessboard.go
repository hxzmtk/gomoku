package chessboard

type Hand uint

const (
	NilHand   Hand = iota //空白
	BlackHand             //黑手
	WhiteHand             //白手
)

func (h Hand) Reverse() Hand {
	switch h {
	case BlackHand:
		return WhiteHand
	case WhiteHand:
		return BlackHand
	default:
		return NilHand
	}
}

type Node interface {
	Go(x, y int, value Hand) error
	IsWin(x, y int) bool
	IsEmpty() bool
	IsFull() bool
	Reset()
	GetState() XYS
	Copy() Node
	Clear(x,y int)
}

var (
	_ Node = &node{}
	_ Node = &nodeArray{}
)

// 棋盘坐标信息
type XY struct {
	X    int  `json:"x"`
	Y    int  `json:"y"`
	Hand Hand `json:"hand"`
}

type XYS []XY

func (xys XYS) NoNilHand() (data []XY) {
	for _, x := range xys {
		if x.Hand != NilHand {
			data = append(data, x)
		}
	}
	return
}
