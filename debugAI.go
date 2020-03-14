package main

import (
	"fmt"
	"strconv"
	"strings"
)

const INCESSANT = Five * 10

type hand uint

type direction uint

const (
	NilHand   hand = iota //空白
	BlackHand             //黑手
	WhiteHand             //白手
)

const (
	left direction = iota
	LeftTop
	LeftBottom
	Top
	Right
	RightTop
	RightBottom
	Bottom
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

func (h hand) Reverse() hand {
	if h == BlackHand {
		return WhiteHand
	} else if h == WhiteHand {
		return BlackHand
	}
	return NilHand
}

type ScoreValue int

const (
	One          ScoreValue = 10
	Two                     = 100
	Three                   = 1000
	Four                    = 10000
	Five                    = 100000
	BlockedOne              = 1
	BlockedTwo              = 10
	BlockedThree            = 100
	BlockedFour             = 1000
)

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
func (g *Grid) Set(row, col int, value hand) {
	offset := g
	offset = g.Offset(row, col)
	if offset == nil {
		return
	} else if offset.value == NilHand {
		offset.value = value
	}
	return
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

func (g *Grid) Print() {
	fmt.Println("当前棋盘布局为:")
	var colNumStr = ""
	col := g.GetColLen()
	fillC := strconv.Itoa(g.GetColLen())
	for i := 1; i <= col; i++ {
		colNumStr += " " + StrLeftFill(len(fillC), i-1)
	}
	fmt.Println(StrLeftFill(len(strconv.Itoa(g.GetRowLen())), ""), strings.TrimLeft(colNumStr, " "))

	for row := 1; row <= g.GetRowLen(); row++ {
		var rowStr = ""
		for col := 1; col <= g.GetColLen(); col++ {
			//rowStr += StrLeftFill(len(strconv.Itoa(g.GetColLen())), "") + g.Offset(row, col).value.Str()
			rowStr += " " + StrLeftFill(len(strconv.Itoa(g.GetColLen())), g.Offset(row, col).value.Str())
		}
		fmt.Println(StrLeftFill(len(strconv.Itoa(g.GetRowLen())), row-1), strings.TrimLeft(rowStr, " "))
	}

	fmt.Println("")
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

func (g *Grid) OffsetXY(x, y int) *Grid {
	return g.Offset(g.XYtoRC(x, y))
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

type Pos struct {
	X int
	Y int
}

func (g *Grid) GetEmptyXY() []Pos {
	var pos []Pos
	for row := 1; row <= g.GetRowLen(); row++ {
		for col := 1; col <= g.GetColLen(); col++ {
			if g.Offset(row, col).value == NilHand {
				pos = append(pos, Pos{
					X: col - 1,
					Y: row - 1,
				})
			}
		}
	}
	return pos
}

func (g *Grid) RCtoXY(row, col int) (x, y int) {
	return col - 1, row - 1
}

func (g *Grid) XYtoRC(x, y int) (row, col int) {
	return y + 1, x + 1
}

func (g *Grid) Copy() *Grid {
	newGrid := InitGrid(g.GetRowLen(), g.GetColLen(), &Grid{})
	for row := 1; row <= g.GetRowLen(); row++ {
		for col := 1; col <= g.GetColLen(); col++ {
			if g.Offset(row, col).value != NilHand {
				newGrid.Offset(row, col).value = g.Offset(row, col).value
			}
		}
	}
	return newGrid
}

func (g *Grid) GetHandXY(h hand) []Pos {
	var pos []Pos
	if h == NilHand {
		return pos
	}
	for row := 1; row <= g.GetRowLen(); row++ {
		for col := 1; col <= g.GetColLen(); col++ {
			if g.Offset(row, col).value == h {
				pos = append(pos, Pos{
					X: col - 1,
					Y: row - 1,
				})
			}
		}
	}
	return pos
}

func (g *Grid) CheckWin(x, y int) bool {

	offset := g.OffsetXY(x, y)
	if offset == nil || offset.value == NilHand {
		return false
	}
	l, _ := g.CountLink(left, x, y)
	lTop, _ := g.CountLink(LeftTop, x, y)
	lBottom, _ := g.CountLink(LeftBottom, x, y)
	t, _ := g.CountLink(Top, x, y)
	b, _ := g.CountLink(Bottom, x, y)
	r, _ := g.CountLink(Right, x, y)
	rTop, _ := g.CountLink(RightTop, x, y)
	rBottom, _ := g.CountLink(RightBottom, x, y)

	//连5
	if l+r >= 4 || lTop+rBottom >= 4 || rTop+lBottom >= 4 || t+b >= 4 {
		return true
	}
	return false
}

func (g *Grid) Win() bool {
	for _, t := range [2]hand{BlackHand, WhiteHand} {
		for _, p := range g.GetHandXY(t) {
			x, y := p.X, p.Y
			l, _ := g.CountLink(left, x, y)
			lTop, _ := g.CountLink(LeftTop, x, y)
			lBottom, _ := g.CountLink(LeftBottom, x, y)
			t, _ := g.CountLink(Top, x, y)
			b, _ := g.CountLink(Bottom, x, y)
			r, _ := g.CountLink(Right, x, y)
			rTop, _ := g.CountLink(RightTop, x, y)
			rBottom, _ := g.CountLink(RightBottom, x, y)

			//连5
			if l+r >= 4 || lTop+rBottom >= 4 || rTop+lBottom >= 4 || t+b >= 4 {
				return true
			}
			return false
		}
	}
	return false
}

//字符串左边填充
func StrLeftFill(s int, value interface{}) string {
	var format = ""
	format = "%" + strconv.Itoa(s) + "v"
	return fmt.Sprintf(format, value)
}

type board struct {
	X     int
	Y     int
	Score int
}

var boards [15][15]int

type AI struct {
	X     int
	Y     int
	Depth int
	h     hand
}

func (ai *AI) negaMax(grid *Grid, h hand, alpha, beta, depth int) int {
	if depth <= 0 || grid.Win() {
		return ai.Evaluate(grid, h)
	}

	for _, p := range grid.GetEmptyXY() {
		//if grid.CheckWin(p.X, p.Y) {
		//	return ai.Evaluate(grid.Copy(), h)
		//}
		if grid.HasNeighbor(p.X, p.Y) {
			g := grid.Copy()
			g.OffsetXY(p.X, p.Y).value = h

			value := -ai.negaMax(g, h.Reverse(), -beta, -alpha, depth-1)
			if value > alpha {
				if depth == ai.Depth {
					ai.X = p.X
					ai.Y = p.Y
					//boards[p.X][p.Y] = value
					//fmt.Println(p.X, p.Y, value, h)
				}

				//boards[p.X][p.Y] = value

				//剪枝
				if value >= beta {
					return beta
				}
				alpha = value
			}
		}
	}
	return alpha
}

//评分
func (ai *AI) Evaluate(g *Grid, h hand) int {
	totalScore := 0
	AIScore := 0
	enemyScore := 0
	computer := ai.h
	target := computer.Reverse()
	for _, t := range [2]hand{computer, target} {
		for _, p := range g.GetHandXY(t) {
			l, lIsSet := g.CountLink(left, p.X, p.Y)
			lTop, lTopIsSet := g.CountLink(LeftTop, p.X, p.Y)
			lBottom, lBottomIsSet := g.CountLink(LeftBottom, p.X, p.Y)
			t, tIsSet := g.CountLink(Top, p.X, p.Y)
			b, bIsSet := g.CountLink(Bottom, p.X, p.Y)
			r, rIsSet := g.CountLink(Right, p.X, p.Y)
			rTop, rTopIsSet := g.CountLink(RightTop, p.X, p.Y)
			rBottom, rBottomIsSet := g.CountLink(RightBottom, p.X, p.Y)

			//连5
			if l+r >= 4 || lTop+rBottom >= 4 || rTop+lBottom >= 4 || t+b >= 4 {
				if t == int(computer) {
					AIScore += Five
				} else {
					enemyScore += Five
				}
			}

			//活四
			if (l+r == 3 && lIsSet && rIsSet) || (lTop+rBottom == 3 && lTopIsSet && rBottomIsSet) ||
				(rTop+lBottom == 3 && rTopIsSet && lBottomIsSet) || (t+b == 3 && tIsSet && bIsSet) {
				if t == int(computer) {
					AIScore += Four
				} else {
					enemyScore += Four
				}
			}

			//活三
			if (l+r == 2 && lIsSet && rIsSet) || (lTop+rBottom == 2 && lTopIsSet && rBottomIsSet) ||
				(rTop+lBottom == 2 && rTopIsSet && lBottomIsSet) || (t+b == 2 && tIsSet && bIsSet) {
				if t == int(computer) {
					AIScore += Three
				} else {
					enemyScore += Three
				}
			}

			//双活三
			if (l == 2 && lIsSet) && ((lTop == 2 && lTopIsSet) || (t == 2 && tIsSet) || (rTop == 2 && rTopIsSet) ||
				(rBottom == 2 && rBottomIsSet) || (b == 2 && bIsSet) || (lBottom == 2 && lBottomIsSet)) {
				if t == int(computer) {
					AIScore += Three * 2
				} else {
					enemyScore += Three * 2
				}
			}

			//活二
			if (l+r == 1 && lIsSet && rIsSet) || (lTop+rBottom == 1 && lTopIsSet && rBottomIsSet) ||
				(rTop+lBottom == 1 && rTopIsSet && lBottomIsSet) || (t+b == 1 && tIsSet && bIsSet) {
				if t == int(computer) {
					AIScore += Two
				} else {
					enemyScore += Two
				}
			}

			//活一
			if (l+r == 0 && lIsSet && rIsSet) || (lTop+rBottom == 0 && lTopIsSet && rBottomIsSet) ||
				(rTop+lBottom == 0 && rTopIsSet && lBottomIsSet) || (t+b == 0 && tIsSet && bIsSet) {
				if t == int(computer) {
					AIScore += int(One)
				} else {
					enemyScore += int(One)
				}
			}

			//死四
			if (l+r == 3 && (!lIsSet || !rIsSet)) || (lTop+lBottom == 3 && (!lTopIsSet || !lBottomIsSet)) ||
				(rTop+rBottom == 3 && (!rBottomIsSet || !rTopIsSet)) || (t+b == 3 && (!tIsSet || !bIsSet)) {
				if t == int(computer) {
					AIScore += BlockedFour
				} else {
					enemyScore += BlockedFour
				}
			}

			//死三
			if (l+r == 2 && (!lIsSet || !rIsSet)) || (lTop+rBottom == 2 && (!lTopIsSet || !rBottomIsSet)) ||
				(rTop+lBottom == 2 && (lBottomIsSet || !rTopIsSet)) || (t+b == 2 && (!tIsSet || !bIsSet)) {
				if t == int(computer) {
					AIScore += BlockedThree
				} else {
					enemyScore += BlockedThree
				}
			}

			//死二
			if (l+r == 1 && (!lIsSet || !rIsSet)) || (lTop+rBottom == 1 && (!lTopIsSet || !rBottomIsSet)) ||
				(rTop+lBottom == 1 && (lBottomIsSet || !rTopIsSet)) || (t+b == 1 && (!tIsSet || !bIsSet)) {
				if t == int(computer) {
					AIScore += BlockedTwo
				} else {
					enemyScore += BlockedTwo
				}
			}

			//死一
			if (l+r == 1 && (!lIsSet || !rIsSet)) || (lTop+rBottom == 1 && (!lTopIsSet || !rBottomIsSet)) ||
				(rTop+lBottom == 1 && (lBottomIsSet || !rTopIsSet)) || (t+b == 1 && (!tIsSet || !bIsSet)) {
				if t == int(computer) {
					AIScore += BlockedOne
				} else {
					enemyScore += BlockedOne
				}
			}
		}
	}
	totalScore = AIScore - enemyScore
	if ai.h == h {
		return totalScore
	}
	return -totalScore
}

//统计每个方向连续有几个棋子,是否还能落子
func (g *Grid) CountLink(d direction, x, y int) (count int, isSet bool) {
	offset := g.OffsetXY(x, y)
	h := g.OffsetXY(x, y).value
	if h == NilHand || offset == nil {
		return
	}

	switch d {
	case left:
		for offset.left != nil && offset.left.value == h {
			offset = offset.left
			count += 1
		}
		if offset != nil && offset.left != nil && offset.left.value == NilHand {
			isSet = true
		}
	case LeftTop:
		for offset.LeftTop() != nil && offset.LeftTop().value == h {
			offset = offset.LeftTop()
			count += 1
		}
		if offset != nil && offset.LeftTop() != nil && offset.LeftTop().value == NilHand {
			isSet = true
		}
	case LeftBottom:
		for offset.LeftBottom() != nil && offset.LeftBottom().value == h {
			offset = offset.LeftBottom()
			count += 1
		}
		if offset != nil && offset.LeftBottom() != nil && offset.LeftBottom().value == NilHand {
			isSet = true
		}
	case Top:
		for offset.top != nil && offset.top.value == h {
			offset = offset.top
			count += 1
		}
		if offset != nil && offset.top != nil && offset.top.value == NilHand {
			isSet = true
		}
	case Right:
		for offset.right != nil && offset.right.value == h {
			offset = offset.right
			count += 1
		}
		if offset != nil && offset.right != nil && offset.right.value == NilHand {
			isSet = true
		}
	case RightTop:
		for offset.RightTop() != nil && offset.RightTop().value == h {
			offset = offset.RightTop()
			count += 1
		}
		if offset != nil && offset.RightTop() != nil && offset.RightTop().value == NilHand {
			isSet = true
		}
	case RightBottom:
		for offset.RightBottom() != nil && offset.RightBottom().value == h {
			offset = offset.RightBottom()
			count += 1
		}
		if offset != nil && offset.RightBottom() != nil && offset.RightBottom().value == NilHand {
			isSet = true
		}
	case Bottom:
		for offset.bottom != nil && offset.bottom.value == h {
			offset = offset.bottom
			count += 1
		}
		if offset != nil && offset.bottom != nil && offset.bottom.value == NilHand {
			isSet = true
		}
	}
	return
}

func (g *Grid) HasNeighbor(x, y int) bool {
	pos := g.OffsetXY(x, y)
	if pos.LeftTop() != nil && pos.LeftTop().value != NilHand {
		return true
	} else if pos.LeftBottom() != nil && pos.LeftBottom().value != NilHand {
		return true
	} else if pos.top != nil && pos.top.value != NilHand {
		return true
	} else if pos.RightTop() != nil && pos.RightTop().value != NilHand {
		return true
	} else if pos.RightBottom() != nil && pos.RightBottom().value != NilHand {
		return true
	} else if pos.bottom != nil && pos.bottom.value != NilHand {
		return true
	}
	return false
}

func main() {
	grid := InitGrid(15, 15, &Grid{})
	grid.Set(7, 7, BlackHand)
	grid.Set(7, 8, WhiteHand)
	grid.Set(8, 7, BlackHand)
	grid.Set(8, 8, WhiteHand)
	grid.Set(9, 7, BlackHand)
	//grid.OffsetXY(5, 5).value = WhiteHand
	//grid.OffsetXY(6, 5).value = BlackHand

	ai := AI{
		Depth: 3,
		h:     WhiteHand,
	}
	//grid.OffsetXY(7, 7).value = BlackHand
	t := WhiteHand
	//i := 0
	grid.Print()
	for !grid.IsFull() {
		ai.negaMax(grid.Copy(), t, -INCESSANT, INCESSANT, ai.Depth)
		grid.OffsetXY(ai.X, ai.Y).value = t
		t = t.Reverse()
		break
		//if i == 100 || grid.CheckWin(ai.X, ai.Y) {
		//	grid.Print()
		//	break
		//}
		//i++
	}
	fmt.Println()
	ai.negaMax(grid.Copy(), t, -INCESSANT, INCESSANT, ai.Depth)
	grid.OffsetXY(ai.X, ai.Y).value = t
	t = t.Reverse()

	ai.negaMax(grid.Copy(), t, -INCESSANT, INCESSANT, ai.Depth)
	grid.OffsetXY(ai.X, ai.Y).value = t
	t = t.Reverse()
	ai.negaMax(grid.Copy(), t, -INCESSANT, INCESSANT, ai.Depth)
	grid.OffsetXY(ai.X, ai.Y).value = t

	//AI(grid.Copy(), t)
	//grid.OffsetXY(x, y).value = t
	grid.Print()
}
