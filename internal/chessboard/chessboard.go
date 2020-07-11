package chessboard

type Hand uint

const (
	NilHand   Hand = iota //空白
	BlackHand             //黑手
	WhiteHand             //白手
)

type Node interface {
	Go(x, y int, value Hand) bool
	IsWin(x, y int) bool
	IsEmpty() bool
	IsFull() bool
	Reset()
	GetState() []XY
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
