package board

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/game"
	"github.com/otrego/clamshell/core/point"
)

// Board Contains the board, capturesStones, and ko
// capturedStones only contains the stones captured
// from the most recent AddStone call. ko contains
// a point that is illegal to recapture due to Ko.
type Board struct {
	board [][]color.Color
	ko    *point.Point
}

// NewBoard creates a new size x size board.
func NewBoard(size int) *Board {
	board := Board{
		make([][]color.Color, size),
		nil,
	}

	for i := 0; i < size; i++ {
		board.board[i] = make([]color.Color, size)
	}
	return &board
}

// AddStone adds a stone to the board.
// returns the captured stones, or err
// if any go (baduk) rules were broken
func (b *Board) AddStone(m *game.Move) ([]*point.Point, error) {
	pt := m.Point()
	var (
		x, y int          = int(pt.X()), int(pt.Y())
		ko   *point.Point = b.ko
	)
	b.ko = nil

	if x >= len(b.board[0]) || y >= len(b.board) ||
		x < 0 || y < 0 {
		return nil, fmt.Errorf("move %v out of bounds for %dx%d board",
			pt, len(b.board[0]), len(b.board))
	}
	if b.board[y][x] != color.Empty {
		return nil, fmt.Errorf("move %v already occupied", pt)
	}

	b.board[y][x] = m.Color()
	capturedStones := b.findCapturedGroups(m)
	if len(capturedStones) == 0 && len(b.capturedStones(m.Point())) != 0 {
		b.board[y][x] = color.Empty
		return nil, fmt.Errorf("move %v is suicidal", pt)
	}
	if len(capturedStones) == 1 {
		b.ko = m.Point()
		if *ko == *(capturedStones[0]) {
			b.board[y][x] = color.Empty
			return nil, fmt.Errorf("%v is an illegal ko move", pt)
		}
	}

	b.removeCapturedStones(capturedStones)
	return capturedStones, nil
}

// findCapturedGroups finds the groups captured by *Move m.
func (b *Board) findCapturedGroups(m *game.Move) []*point.Point {
	pt := m.Point()

	points := make([]*point.Point, 4)
	points[0] = point.New(pt.X()+1, pt.Y())
	points[1] = point.New(pt.X()-1, pt.Y())
	points[2] = point.New(pt.X(), pt.Y()+1)
	points[3] = point.New(pt.X(), pt.Y()-1)

	capturedStones := make([]*point.Point, 0)
	for _, point := range points {
		capturedStones = append(capturedStones, b.capturedStones(point)...)
	}
	return capturedStones
}

// removeCapturedStones removes the captured stones from
// the board.
func (b *Board) removeCapturedStones(capturedStones []*point.Point) {
	for _, point := range capturedStones {
		b.board[point.Y()][point.X()] = color.Empty
	}
}

// CapturedStones returns the captured stones in group containing Point pt.
// returns nil if no stones were captured
func (b *Board) capturedStones(pt *point.Point) []*point.Point {
	expanded := make(map[point.Point]bool)

	// current group color
	c := b.board[pt.Y()][pt.X()]

	queue := list.New()
	queue.PushBack(pt)
	for queue.Len() > 0 {
		e := queue.Front()
		queue.Remove(e)
		pt1, ok := e.Value.(*point.Point)
		if !ok {
			panic("e.Value was not of type point.Point")
		}

		x, y := int(pt1.X()), int(pt1.Y())
		if x >= len(b.board[0]) || y >= len(b.board) ||
			x < 0 || y < 0 {
			continue
		} else if b.board[y][x] == color.Empty {
			// Liberty has been found, no need to continue search
			return nil
		} else if b.board[y][x] == c && !expanded[*pt1] {
			expanded[*pt1] = true
			queue.PushBack(point.New(pt1.X()+1, pt1.Y())) // enqueue right
			queue.PushBack(point.New(pt1.X()-1, pt1.Y())) // enqueue left
			queue.PushBack(point.New(pt1.X(), pt1.Y()+1)) // enqueue down
			queue.PushBack(point.New(pt1.X(), pt1.Y()-1)) // enqueue up
		}
	}

	// The stones that were captured
	stoneGroup := make([]*point.Point, len(expanded))
	i := 0
	for key := range expanded {
		stoneGroup[i] = point.New(key.X(), key.Y())
		i++
	}
	return stoneGroup
}

// String returns a string representation of this board.
// For example:
//
//    b.Board {{B, W, B,  },
//             {W,  , B, B},
//             { ,  , W,  },
//             {B,  , W,  }}
//
//    Becomes  [B W B .]
//             [W . B B]
//             [. . W .]
//             [B . W .]
func (b *Board) String() string {
	var sb strings.Builder
	for i := 0; i < len(b.board); i++ {
		// To increase useability of this String function,
		// color.Empty is converted from "" to ".".
		str := make([]string, len(b.board[0]))
		for j := 0; j < len(b.board[0]); j++ {
			if b.board[i][j] == color.Empty {
				str[j] = "."
			} else {
				str[j] = string(b.board[i][j])
			}
		}
		sb.WriteString(fmt.Sprintf("%v\n", str))
	}
	return strings.TrimSpace(sb.String())
}
