package board

import (
	"container/list"
	"errors"
	"fmt"
	"strings"

	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/move"
	"github.com/otrego/clamshell/go/point"
)

var (
	InvalidBoardState = errors.New("invalid board state")
	IllegalMove       = errors.New("illegal move")
)

// Board Contains the board, capturesStones, and ko
// ko contains a point that is illegal to recapture due to Ko.
type Board struct {
	// The board, arranged in rows (rather than columns).
	board [][]color.Color
	ko    *point.Point
}

// New creates a new size x size board.
func New(size int) *Board {
	board := Board{
		make([][]color.Color, size),
		nil,
	}

	for i := 0; i < size; i++ {
		board.board[i] = make([]color.Color, size)
	}
	return &board
}

// PlaceStone adds a stone to the board and removes captured stones (if any).
// returns the captured stones, or err if any Go (baduk) rules were broken
func (b *Board) PlaceStone(m *move.Move) (move.List, error) {
	if !b.inBounds(m.Point()) {
		return nil, fmt.Errorf("%w: move %v out of bounds for %dx%d board",
			IllegalMove, m.Point(), len(b.board[0]), len(b.board))
	}
	if b.colorAt(m.Point()) != color.Empty {
		return nil, fmt.Errorf("%w: move %v already occupied", IllegalMove, m.Point())
	}

	b.setColor(m)

	capturedStones := b.findCapturedGroups(m)
	if len(capturedStones) == 0 && len(b.capturedStones(m.Point())) != 0 {
		b.setColor(move.New(color.Empty, m.Point()))
		return nil, fmt.Errorf("%w: move %v is suicidal", IllegalMove, m.Point())
	}
	if len(capturedStones) == 1 {
		if b.ko != nil && *(b.ko) == *(m.Point()) {
			b.setColor(move.New(color.Empty, m.Point()))
			return nil, fmt.Errorf("%w: %v is an illegal ko move", IllegalMove, m.Point())
		}
		b.ko = capturedStones[0]
	} else {
		b.ko = nil
	}

	// convert the captured stones into Move objects for convience.
	var captured move.List
	opp := m.Color().Opposite()
	for _, pt := range capturedStones {
		captured = append(captured, move.New(opp, pt))
	}
	captured.Sort()

	b.removeCapturedStones(capturedStones)
	return captured, nil
}

// findCapturedGroups returns the groups captured by *Move m.
func (b *Board) findCapturedGroups(m *move.Move) []*point.Point {
	pt := m.Point()

	points := b.getNeighbors(pt)
	capturedStones := make([]*point.Point, 0)
	for _, point := range points {
		if b.inBounds(point) {
			capturedStones = append(capturedStones, b.capturedStones(point)...)
		}
	}
	return capturedStones
}

// removeCapturedStones removes the captured stones from
// the board.
func (b *Board) removeCapturedStones(capturedStones []*point.Point) {
	for _, point := range capturedStones {
		b.setColor(move.New(color.Empty, point))
	}
}

// capturedStones returns the captured stones in group containing Point pt.
// returns nil if no stones were captured.
func (b *Board) capturedStones(pt *point.Point) []*point.Point {
	stoneGroup, captured := b.getStoneGroup(pt)
	if captured {
		return stoneGroup
	}
	return nil
}

// getStoneGroup returns all the points in a stone group containing point pt
// and true if the group is captured
func (b *Board) getStoneGroup(pt *point.Point) ([]*point.Point, bool) {
	expanded := make(map[point.Point]bool)
	captured := true
	// current group color
	c := b.colorAt(pt)

	queue := list.New()
	queue.PushBack(pt)
	for queue.Len() > 0 {
		e := queue.Front()
		queue.Remove(e)
		pt1, ok := e.Value.(*point.Point)
		if !ok {
			panic("e.Value was not of type point.Point")
		}

		if !b.inBounds(pt1) {
			continue
		} else if b.colorAt(pt1) == color.Empty {
			captured = false
		} else if b.colorAt(pt1) == c && !expanded[*pt1] {
			expanded[*pt1] = true
			points := b.getNeighbors(pt1)
			for _, point := range points {
				queue.PushBack(point)
			}
		}
	}

	// The stones in this group
	stoneGroup := make([]*point.Point, len(expanded))
	i := 0
	for key := range expanded {
		stoneGroup[i] = point.New(key.X(), key.Y())
		i++
	}
	return stoneGroup, captured
}

// inBounds returns true if x and y are in bounds
// on the board, false otherwise.
func (b *Board) inBounds(pt *point.Point) bool {
	var x, y int = pt.X(), pt.Y()
	return x < len(b.board[0]) && y < len(b.board) &&
		x >= 0 && y >= 0
}

// colorAt returns the color at point pt.
func (b *Board) colorAt(pt *point.Point) color.Color {
	var x, y int = pt.X(), pt.Y()
	return b.board[y][x]
}

// setColor sets the color m.Color at point m.Point.
func (b *Board) setColor(m *move.Move) {
	var x, y int = m.Point().X(), m.Point().Y()
	b.board[y][x] = m.Color()
}

// getNeighbors returns a list of points neighboring point pt.
// Neighboring points could be out of bounds.
func (b *Board) getNeighbors(pt *point.Point) []*point.Point {
	points := make([]*point.Point, 4)
	points[0] = point.New(pt.X()+1, pt.Y())
	points[1] = point.New(pt.X()-1, pt.Y())
	points[2] = point.New(pt.X(), pt.Y()+1)
	points[3] = point.New(pt.X(), pt.Y()-1)
	return points
}

// SetPlacements force-places moves on the go-board, without performing capture
// logic. If an illegal board position results, return an error.
func (b *Board) SetPlacements(ml move.List) error {

	for _, m := range ml {
		b.setColor(m)
	}

	// Validate we have a valid board position -- i.e., one
	// without captures lying on the board.
	explored := make(map[point.Point]bool)
	for _, m := range ml {
		pt := m.Point()

		if !explored[*pt] {
			stoneGroup, captured := b.getStoneGroup(pt)
			if captured {
				return fmt.Errorf("%w: stones at points %v are captured", InvalidBoardState, stoneGroup)
			}
			for _, point := range stoneGroup {
				explored[*point] = true
			}
		}
	}
	return nil
}

// Ko returns the ko point.
func (b *Board) Ko() *point.Point {
	return b.ko
}

// Clone makes a board copy.
func (b *Board) Clone() *Board {
	newb := &Board{
		ko:    b.ko,
		board: make([][]color.Color, len(b.board)),
	}
	for i, row := range b.board {
		newRow := make([]color.Color, len(row))
		for j, col := range row {
			newRow[j] = col
		}
		newb.board[i] = newRow
	}
	return newb
}

// StoneState returns an array of all the current stone positions.
func (b *Board) StoneState() move.List {
	var moves move.List
	for i := 0; i < len(b.board); i++ {
		for j := 0; j < len(b.board[0]); j++ {
			if b.board[i][j] != color.Empty {
				moves = append(moves,
					move.New(b.board[i][j], point.New(j, i)))
			}
		}
	}
	return moves
}

// FullBoardState returns the full board state.
func (b *Board) FullBoardState() [][]color.Color {
	out := make([][]color.Color, len(b.board))
	for i, row := range b.board {
		out[i] = make([]color.Color, len(row))
		for j, col := range row {
			out[i][j] = col
		}
	}
	return out
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
			if b.ko != nil && b.ko.X() == j && b.ko.Y() == i {
				str[j] = "*"
			} else if b.board[i][j] == color.Empty {
				str[j] = "."
			} else {
				str[j] = string(b.board[i][j])
			}
		}
		sb.WriteString(fmt.Sprintf("%v\n", str))
	}
	return strings.TrimSpace(sb.String())
}
