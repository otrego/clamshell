package move

import (
	"fmt"

	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/point"
)

// A Move is location + a color. A pass is represented by a Move with defined
// color but no point specified.
type Move struct {
	color color.Color
	point *point.Point
}

// New creates a new Move.
func New(col color.Color, pt *point.Point) *Move {
	return &Move{color: col, point: pt}
}

// NewPass creates a a new pass Move.
func NewPass(col color.Color) *Move {
	return &Move{color: col, point: nil}
}

// Color returns the color.
func (m *Move) Color() color.Color {
	return m.color
}

// Point returns the point.
func (m *Move) Point() *point.Point {
	return m.point
}

// String returns the string value for a Move.
func (m *Move) String() string {
	return fmt.Sprintf("{%v, %v}", m.color, m.point)
}

// GoString returns the string value.
func (m *Move) GoString() string {
	return fmt.Sprintf("{Color:%v, Point:%v}", m.color, m.point)
}

// IsPass indicates whether this is a 'pass' move (i.e., there is a player but
// no point).
func (m *Move) IsPass() bool {
	return m.point == nil
}

// FromSGFPoint converts from an SGF point of the form "ab" to a point
// object, such as {0,1}.
func FromSGFPoint(col color.Color, sgfPt string) (*Move, error) {
	if sgfPt == "" {
		// This is valid. This is a 'Pass' Move.
		return &Move{color: col}, nil
	}
	pt, err := point.NewFromSGF(sgfPt)
	if err != nil {
		return nil, err
	}
	return &Move{color: col, point: pt}, nil
}

// ListFromSGFPoints a move list of the form "ab", "bc" to a moves of the form
// {0,1}, {0,2}. Note that pass-moves are not allowed in move-lists.
func ListFromSGFPoints(col color.Color, sgfPts []string) ([]*Move, error) {
	var moves []*Move
	for _, sgfPt := range sgfPts {
		pt, err := point.NewFromSGF(sgfPt)
		if err != nil {
			return nil, err
		}
		moves = append(moves, &Move{
			color: col,
			point: pt,
		})
	}
	return moves, nil
}
