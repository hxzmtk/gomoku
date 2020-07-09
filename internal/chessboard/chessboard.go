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

func (n *node) rightBottom() *node {
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

//检查是否已分出胜负
func (n *node) IsWin(x, y int) bool {
	offset := n.get(x, y)
	h := offset.value
	if h == NilHand {
		return false
	}

	//检查行
	count := 1
	left := offset.left
	right := offset.right
	for {
		if left == nil {
			break
		}
		if left.value == h {
			count++
		} else {
			break
		}
		left = left.left
	}
	for {
		if right == nil {
			break
		}
		if right.value == h {
			count++
		} else {
			break
		}
		right = right.right
	}
	if count >= 5 {
		return true
	}

	//检查列
	count = 1
	top := offset.top
	bottom := offset.bottom
	for {
		if top == nil {
			break
		}
		if top.value == h {
			count++
		} else {
			break
		}
		top = top.top
	}
	for {
		if bottom == nil {
			break
		}
		if bottom.value == h {
			count++
		} else {
			break
		}
		bottom = bottom.bottom
	}
	if count >= 5 {
		return true
	}

	//检查左斜边
	count = 1
	leftTop := offset.leftTop()
	rightBottom := offset.rightBottom()
	for {
		if leftTop == nil {
			break
		}
		if leftTop.value == h {
			count++
		} else {
			break
		}
		leftTop = leftTop.leftTop()
	}
	for {
		if rightBottom == nil {
			break
		}
		if rightBottom.value == h {
			count++
		} else {
			break
		}
		rightBottom = rightBottom.rightBottom()
	}
	if count >= 5 {
		return true
	}

	//检查右斜边
	count = 1
	rightTop := offset.rightTop()
	leftBottom := offset.leftBottom()
	for {
		if rightTop == nil {
			break
		}
		if rightTop.value == h {
			count++
		} else {
			break
		}
		rightTop = rightTop.rightTop()
	}
	for {
		if leftBottom == nil {
			break
		}
		if leftBottom.value == h {
			count++
		} else {
			break
		}
		leftBottom = leftBottom.leftBottom()
	}
	if count >= 5 {
		return true
	}
	if n.IsFull() {
		return true
	} else {
		return false
	}
}

//棋盘是否已满？
func (n *node) IsFull() bool {
	for x := 0; x < n.getWidth(); x++ {
		for y := 0; y < n.getHeight(); y++ {
			offset := n.get(x, y)
			if offset != nil && offset.value == NilHand {
				return false
			}
		}
	}
	return true
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
