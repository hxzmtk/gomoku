package chessboard

/*
二维数组实现的棋盘
*/

type nodeArray [][]Hand

func (n *nodeArray) Go(x, y int, value Hand) bool {
	return true
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

func NewChessboardWithArray(size int) *nodeArray {
	array := make(nodeArray, size)
	for i, _ := range array {
		array[i] = make([]Hand, size)
	}
	return &array
}
