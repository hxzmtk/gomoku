package gomoku

import (
	"fmt"
)

const INCESSANT = Five * 10

type direction uint

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

func (h Hand) Reverse() Hand {
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

func (g *Grid) OffsetXY(x, y int) *Grid {
	return g.Offset(g.XYtoRC(x, y))
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

func (g *Grid) GetHandXY(h Hand) []Pos {
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
		if offset.value == WhiteHand {
			fmt.Println("白手赢", fmt.Sprintf("最后一个落子点为(x:%d,y:%d)", x, y))
		} else if offset.value == BlackHand {
			fmt.Println("黑手赢", fmt.Sprintf("最后一个落子点为(x:%d,y:%d)", x, y))
		}
		return true
	}
	return false
}

func (g *Grid) Win() bool {
	for _, t := range [2]Hand{BlackHand, WhiteHand} {
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
	H     Hand
}

func (ai *AI) NegaMax(grid *Grid, h Hand, alpha, beta, depth int) int {
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

			value := -ai.NegaMax(g, h.Reverse(), -beta, -alpha, depth-1)
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
func (ai *AI) Evaluate(g *Grid, h Hand) int {
	totalScore := 0
	AIScore := 0
	enemyScore := 0
	computer := ai.H
	target := computer.Reverse()
	for _, t := range [2]Hand{computer, target} {
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
	if ai.H == h {
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
