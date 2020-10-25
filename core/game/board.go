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
	if x >= len(board[0]) || y >= len(board) {
		return fmt.Errorf("point (%d, %d) out of bound for %dx%d board",
			x, y, len(board[0]), len(board))
	}
	board[y][x] = c
	return nil
}

// FindCaptures finds all the captured stones on the board
// and calls RemoveCaptures. TODO
func (board Board) FindCaptures() {
	var expanded = make([][]bool, len(board))
	for i := 0; i < len(expanded); i++ {
		expanded[i] = make([]bool, len(board[0]))
	}
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[0]); x++ {
			// skip empty or expanded nodes.
			if board[y][x] != "" || !expanded[y][x] {
				//current group color
				c := board[y][x]
				//FIXE ME maybe move the contents of if statement to its own function
				//becomes true only if a liberty is found for this group
				liberty := false
				queue := list.New()
				queue.PushBack(point.New(int64(x), int64(y)))
				for queue.Len() > 0 {
					e := queue.Front()
					queue.Remove(e)
					pt, ok := e.Value.(point.Point)
					if !ok {
						panic("e.Value was not of type point.Point")
					}
					x1, y1 := pt.X(), pt.Y()
					if board[y1][x1] == "" {
						liberty = true
					} else if board[y1][x1] == c {
						expanded[y1][x1] = true
						if x1+1 < int64(len(board[0])) {
							// enqueue right
							queue.PushBack(point.New(pt.X()+int64(1), pt.Y()))
						}
						if x1-1 >= 0 {
							// enqueue left
							queue.PushBack(point.New(pt.X()+int64(1), pt.Y()))
						}
						if y1+1 < int64(len(board)) {
							// enqueue down
							queue.PushBack(point.New(pt.X(), pt.Y()+int64(1)))
						}
						if y1-1 < 0 {
							// enqueue up
							queue.PushBack(point.New(pt.X(), pt.Y()-int64(1)))
						}
					}
				}
				if !liberty {
					//TODO remove that group!
				}
			}
		}
	}

}

// String returns a string representation of this board.
func (board Board) String() string {
	var sb strings.Builder
	for i := 0; i < len(board); i++ {
		// To increase useability of this String function,
		// color.Empty is converted from "" to ".".
		str := make([]string, len(board[0]))
		for j := 0; j < len(board[0]); j++ {
			if board[i][j] == "" {
				str[j] = "."
			} else {
				str[j] = string(board[i][j])
			}
		}
		sb.WriteString(fmt.Sprintf("%v\n", str))
	}
	return strings.TrimSpace(sb.String())
}
