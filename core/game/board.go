package game

//plan for capture is to bfs every group we find for any liberties.
//make sure to mark expanded nodes!!
import (
	"container/list"
	"fmt"
	"strings"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/point"
)

// Board TODO
type Board [][]color.Color

// NewBoard creates a new size x size board.
func NewBoard(size int) *Board {
	var board Board = make([][]color.Color, size)
	for i := 0; i < size; i++ {
		board[i] = make([]color.Color, size)
	}
	return &board
}

// AddStone adds a stone to the board. TODO
func (board Board) AddStone(pt *point.Point, c color.Color) error {
	var x, y int = int(pt.X()), int(pt.Y())
	switch {
	case x >= len(board[0]) || y >= len(board) ||
		x < 0 || y < 0:
		return fmt.Errorf("move %v out of bound for %dx%d board",
			pt, len(board[0]), len(board))
	case board[y][x] != color.Empty:
		return fmt.Errorf("move %v already occupied", pt)
	case false:
		return fmt.Errorf("move %v is suicidal", pt)
	case false:
		return fmt.Errorf("%v is an illigel ko move", pt)
	}

	board[y][x] = c
	return nil
}

// FindEnemyGroups looks for all the enemy groups next to the stone
// that was placed at Point pt.
func (board Board) FindEnemyGroups(pt *point.Point, c color.Color) {
	points := make([]*point.Point, 4)
	points[0] = point.New(pt.X()+1, pt.Y())
	points[1] = point.New(pt.X()-1, pt.Y())
	points[2] = point.New(pt.X(), pt.Y()+1)
	points[3] = point.New(pt.X(), pt.Y()-1)
	for _, point := range points {
		captured, capPoints := board.IsCaptured(point)
		if captured {
			for _, capPoint := range capPoints {
				board[capPoint.Y()][capPoint.X()] = color.Empty
			}
		}
	}
}

// IsCaptured looks at the group at Point pt and returns true if that
// group is captured, and false otherwise.
func (board Board) IsCaptured(pt *point.Point) (bool, []*point.Point) {
	expanded := make(map[point.Point]bool)
	// current group color
	c := board[pt.Y()][pt.X()]
	// becomes true only if a liberty is found for this group
	captured := true
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
		if x >= len(board[0]) || y >= len(board) ||
			x < 0 || y < 0 {
			continue
		} else if board[y][x] == color.Empty {
			captured = false
			var stoneGroup []*point.Point = nil
			return captured, stoneGroup
		} else if board[y][x] == c && !expanded[*pt1] {
			expanded[*pt1] = true
			// enqueue right
			queue.PushBack(point.New(pt1.X()+1, pt1.Y()))
			// enqueue left
			queue.PushBack(point.New(pt1.X()-1, pt1.Y()))
			// enqueue down
			queue.PushBack(point.New(pt1.X(), pt1.Y()+1))
			// enqueue up
			queue.PushBack(point.New(pt1.X(), pt1.Y()-1))
		}
	}
	stoneGroup := make([]*point.Point, len(expanded)) //the stones that were captured
	i := 0
	for key := range expanded {
		stoneGroup[i] = point.New(key.X(), key.Y())
		i++
	}
	return captured, stoneGroup
}

// String returns a string representation of this board.
func (board Board) String() string {
	var sb strings.Builder
	for i := 0; i < len(board); i++ {
		// To increase useability of this String function,
		// color.Empty is converted from "" to ".".
		str := make([]string, len(board[0]))
		for j := 0; j < len(board[0]); j++ {
			if board[i][j] == color.Empty {
				str[j] = "."
			} else {
				str[j] = string(board[i][j])
			}
		}
		sb.WriteString(fmt.Sprintf("%v\n", str))
	}
	return strings.TrimSpace(sb.String())
}
