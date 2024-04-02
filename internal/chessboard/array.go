package chessboard

/*
二维数组实现的棋盘
*/

type nodeArray [][]Hand

func (n *nodeArray) Go(x, y int, value Hand) error {
	return nil
}

func (n *nodeArray) IsWin(x, y int) bool {
	return true
}

func (n *nodeArray) IsEmpty() bool {
	return true
}

func (n *nodeArray) IsFull() bool {
	return true
}

func (n *nodeArray) Reset() {
}

func (n *nodeArray) GetState() (xy XYS) {
	return
}

func (n *nodeArray) Copy() Node {
	return nil
}

func (n *nodeArray) Clear(x, y int) {
	return
}

func NewChessboardWithArray(size int) *nodeArray {
	array := make(nodeArray, size)
	for i := range array {
		array[i] = make([]Hand, size)
	}
	return &array
}
