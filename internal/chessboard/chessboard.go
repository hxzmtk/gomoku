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
}
