package gomoku

import (
	"fmt"
	"strconv"
	"strings"
)

type hand uint

const (
	NilHand   hand = iota //空白
	BlackHand             //黑手
	WhiteHand             //白手
)

func (h hand) Str() string {
	switch h {
	case NilHand:
		return "."
	case BlackHand:
		return "X"
	case WhiteHand:
		return "O"
	default:
		return "."
	}
}

type Grid struct {
	value  hand
	left   *Grid
	right  *Grid
	top    *Grid
	bottom *Grid
}

func (g *Grid) LeftTop() *Grid {
	if g.left != nil {
		return g.left.top
	}
	return nil
}

func (g *Grid) LeftBottom() *Grid {
	if g.left != nil {
		return g.left.bottom
	}
	return nil
}

func (g *Grid) RightTop() *Grid {
	if g.right != nil {
		return g.right.top
	}
	return nil
}

func (g *Grid) RightBottom() *Grid {
	if g.right != nil {
		return g.right.bottom
	}
	return nil
}

//设置row,col坐标的值, 即 落棋子
func (g *Grid) Set(row, col int, value hand) bool {
	offset := g
	offset = g.Offset(row, col)
	if offset == nil {
		return false
	} else if offset.value == NilHand {
		offset.value = value
	} else if offset.value != NilHand {
		return false
	}
	return true
}

func (g *Grid) SetByXY(x, y int, value hand) bool {
	row, col := y+1, x+1
	return g.Set(row, col, value)
}

//获取向右偏移x位的指针
func (g *Grid) RightOffset(col int) *Grid {
	var tmp *Grid
	tmp = g
	if g == nil {
		return nil
	}
	for i := 1; i <= col; i++ {
		if i == col && tmp != nil {
			return tmp
		}
		if tmp == nil {
			return nil
		}
		if tmp.right != nil {
			tmp = tmp.right
		}
	}
	return nil
}

//获取向下偏移y位的指针
func (g *Grid) BottomOffset(row int) *Grid {
	var tmp *Grid
	tmp = g
	if g == nil {
		return nil
	}
	for i := 1; i <= row; i++ {
		if i == row && tmp != nil {
			return tmp
		}
		if tmp == nil {
			return nil
		}
		if tmp.bottom != nil {
			tmp = tmp.bottom
		}
	}
	return nil
}

//获取该表格有几行
func (g *Grid) GetRowLen() int {
	var row int
	tmp := g
	for tmp != nil {
		row++
		tmp = tmp.bottom
	}
	return row
}

//获取该表格有几列
func (g *Grid) GetColLen() int {
	var col int
	tmp := g
	for tmp != nil {
		col++
		tmp = tmp.right
	}
	return col
}

/*
获取坐标处的指针
*/
func (g *Grid) Offset(row, col int) *Grid {
	offset := g
	if g == nil {
		return nil
	}
	var i, j int
	i, j = 1, 1
	for i <= row {
		if row == i {
			for j <= col {
				if col == j {
					return offset
				}

				if offset == nil {
					return nil
				}
				offset = offset.right
				j++
			}
		}
		if offset == nil {
			return nil
		}
		offset = offset.bottom
		i++
	}
	return nil
}

//获取某列最后一行的棋子的指针
func (g *Grid) GetLastRow(col int) *Grid {
	offset := g.Offset(1, col)
	for row := 1; row <= g.GetRowLen(); row++ {
		if row == g.GetRowLen() {
			return offset
		}
		offset = offset.bottom
	}
	return nil
}

//获取某一行的最后一列没有落棋子的指针
func (g *Grid) GetEmptyLastRow(col int) *Grid {
	if col <= 0 || col > g.GetColLen() {
		return nil
	}
	last := g.GetLastRow(col)
	if last == nil {
		return nil
	}
	for row := 1; row <= g.GetRowLen(); row++ {
		if last.value == NilHand {
			return last
		}
		last = last.top
		if last == nil {
			return nil
		}
	}
	return nil
}

//棋盘是否已满？
func (g *Grid) IsFull() bool {
	for row := 1; row <= g.GetRowLen(); row++ {
		for col := 1; col <= g.GetColLen(); col++ {
			if g.Offset(row, col).value == NilHand {
				return false
			}
		}
	}
	return true
}

//棋盘是否为空
func (g *Grid) IsEmpty() bool {
	for row := 1; row <= g.GetRowLen(); row++ {
		for col := 1; col <= g.GetColLen(); col++ {
			if g.Offset(row, col).value != NilHand {
				return false
			}
		}
	}
	return true
}

//统计黑手,白手的棋子数量
func (g *Grid) Count() (black, white int) {
	for row := 1; row <= g.GetRowLen(); row++ {
		for col := 1; col <= g.GetColLen(); col++ {
			if g.Offset(row, col).value == BlackHand {
				black++
			} else if g.Offset(row, col).value == WhiteHand {
				white++
			}
		}
	}
	return
}

//检查是否已分出胜负
func (g *Grid) IsWin(row, col int) bool {
	offset := g.Offset(row, col)
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
		switch h {
		case WhiteHand:
			fmt.Println("白手赢", fmt.Sprintf("最后一个落子点为(row:%d,col:%d)", row, col))
		case BlackHand:
			fmt.Println("黑手赢", fmt.Sprintf("最后一个落子点为(row:%d,col:%d)", row, col))
		}
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
		switch h {
		case WhiteHand:
			fmt.Println("白手赢", fmt.Sprintf("最后一个落子点为(row:%d,col:%d)", row, col))
		case BlackHand:
			fmt.Println("黑手赢", fmt.Sprintf("最后一个落子点为(row:%d,col:%d)", row, col))
		}
		return true
	}

	//检查左斜边
	count = 1
	leftTop := offset.LeftTop()
	rightBottom := offset.RightBottom()
	for {
		if leftTop == nil {
			break
		}
		if leftTop.value == h {
			count++
		} else {
			break
		}
		leftTop = leftTop.LeftTop()
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
		rightBottom = rightBottom.RightBottom()
	}
	if count >= 5 {
		switch h {
		case WhiteHand:
			fmt.Println("白手赢", fmt.Sprintf("最后一个落子点为(row:%d,col:%d)", row, col))
		case BlackHand:
			fmt.Println("黑手赢", fmt.Sprintf("最后一个落子点为(row:%d,col:%d)", row, col))
		}
		return true
	}

	//检查右斜边
	count = 1
	rightTop := offset.RightTop()
	leftBottom := offset.LeftBottom()
	for {
		if rightTop == nil {
			break
		}
		if rightTop.value == h {
			count++
		} else {
			break
		}
		rightTop = rightTop.RightTop()
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
		leftBottom = leftBottom.LeftBottom()
	}
	if count >= 5 {
		switch h {
		case WhiteHand:
			fmt.Println("白手赢", fmt.Sprintf("最后一个落子点为(row:%d,col:%d)", row, col))
		case BlackHand:
			fmt.Println("黑手赢", fmt.Sprintf("最后一个落子点为(row:%d,col:%d)", row, col))
		}
		return true
	}
	if g.IsFull() {
		fmt.Println("平手")
		return true
	} else {
		return false
	}
}

func (g *Grid) IsWinByXY(x, y int) bool {
	row, col := y+1, x+1
	return g.IsWin(row, col)
}

//重置棋盘
func (g *Grid) Reset() {
	for row := 1; row <= g.GetRowLen(); row++ {
		for col := 1; col <= g.GetColLen(); col++ {
			g.Offset(row, col).value = NilHand
		}
	}
}

/*
约定: row,col的起始值为1
*/
func InitGrid(row, col int, head *Grid) *Grid {
	l := &Grid{}

	l = head
	c := col
	for c > 1 {
		var tmp Grid
		tmp.left = l
		l.right = &tmp
		l = &tmp
		c--
	}

	l = head
	r := row
	for r > 1 {
		var tmp Grid
		tmp.top = l
		l.bottom = &tmp
		l = &tmp
		r--
	}

	for i := 2; i <= row; i++ {
		for j := 2; j <= col; j++ {
			top := head.Offset(i-1, j)
			left := head.Offset(i, j-1)
			var tmp Grid
			tmp.top = top
			tmp.left = left

			top.bottom = &tmp
			left.right = &tmp
		}
	}
	return head
}

func (g *Grid) Print() {
	fmt.Println("当前棋盘布局为:")
	var colNumStr = ""
	col := g.GetColLen()
	fillC := strconv.Itoa(g.GetColLen())
	for i := 1; i <= col; i++ {
		colNumStr += " " + StrLeftFill(len(fillC), i)
	}
	fmt.Println(StrLeftFill(len(strconv.Itoa(g.GetRowLen())), ""), strings.TrimLeft(colNumStr, " "))

	for row := 1; row <= g.GetRowLen(); row++ {
		var rowStr = ""
		for col := 1; col <= g.GetColLen(); col++ {
			//rowStr += StrLeftFill(len(strconv.Itoa(g.GetColLen())), "") + g.Offset(row, col).value.Str()
			rowStr += " " + StrLeftFill(len(strconv.Itoa(g.GetColLen())), g.Offset(row, col).value.Str())
		}
		fmt.Println(StrLeftFill(len(strconv.Itoa(g.GetRowLen())), row), strings.TrimLeft(rowStr, " "))
	}

	fmt.Println("")
}

//字符串左边填充
func StrLeftFill(s int, value interface{}) string {
	var format = ""
	format = "%" + strconv.Itoa(s) + "v"
	return fmt.Sprintf(format, value)
}
