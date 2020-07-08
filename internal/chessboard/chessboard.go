package chessboard

type Hand uint

const (
	NilHand   Hand = iota //空白
	BlackHand             //黑手
	WhiteHand             //白手
)

type Node interface {
	Go(x, y int, value Hand) bool
	IsEmpty() bool
	Reset()
}

type node struct {
	value  Hand
	left   *node
	right  *node
	top    *node
	bottom *node
}

func (n *node) leftTop() *node {
	if n.left != nil {
		return n.left.top
	}
	return nil
}

func (n *node) leftBottom() *node {
	if n.left != nil {
		return n.left.bottom
	}
	return nil
}

func (n *node) rightTop() *node {
	if n.right != nil {
		return n.right.top
	}
	return nil
}

func (n *node) lightBottom() *node {
	if n.right != nil {
		return n.right.bottom
	}
	return nil
}

/*
落子
横坐标为x，纵坐标为y
	x
 ----------->
|
|
| y
|
V
*/
func (n *node) Go(x, y int, value Hand) bool {
	offset := n.get(x, y)
	if offset != nil && offset.value == NilHand {
		offset.value = value
		return true
	}
	return false
}

// 根据坐标获取节点
func (n *node) get(x, y int) *node {
	offset := n
	i, j := 0, 0
	for i <= x {
		if i == x {
			for j <= y {
				if j == y {
					return offset
				}
				if offset == nil {
					return nil
				}
				offset = offset.bottom
				j++
			}
		}
		if offset == nil {
			return nil
		}
		offset = offset.right
		i++
	}
	return nil
}
func (n *node) getWidth() int {
	offset := n
	width := 0
	for offset != nil {
		width++
		offset = offset.right
	}
	return width
}

func (n *node) getHeight() int {
	height := 0
	offset := n
	for offset != nil {
		height++
		offset = offset.bottom
	}
	return height
}

// size为棋盘有几格
func NewChessboard(size int) *node {
	root := new(node)

	offsetX := root
	offsetY := root
	x, y := size, size

	// 连接网格第一行的左右节点
	for x > 1 {
		tmp := &node{
			left: offsetX,
		}
		offsetX.right = tmp
		offsetX = tmp
		x--
	}

	// 连接网格第一列的上下节点
	for y > 1 {
		tmp := &node{
			top: offsetY,
		}
		offsetY.bottom = tmp
		offsetY = tmp
		y--
	}

	// 构造剩余元素，并把left,right,top,bottom都连接起来
	for i := 1; i < size; i++ {
		for j := 1; j < size; j++ {
			n := new(node)           //新建当前节点
			top := root.get(i, j-1)  //当前节点的top节点
			left := root.get(i-1, j) //当前节点的left节点

			n.top = top
			n.left = left
			top.bottom = n
			left.right = n
		}
	}
	return root
}

func (n *node) IsEmpty() bool {
	for x := 0; x < n.getWidth(); x++ {
		for y := 0; y < n.getHeight(); y++ {
			offset := n.get(x, y)
			if offset != nil && offset.value != NilHand {
				return false
			}
		}
	}
	return true
}

func (n *node) Reset() {
	for x := 0; x < n.getWidth(); x++ {
		for y := 0; y < n.getHeight(); y++ {
			offset := n.get(x, y)
			offset.value = NilHand
		}
	}
}
