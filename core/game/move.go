package game

import (
	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/point"
)

// A Move is location + a color. A pass is represented by a Move with defined
// color but no point specified.
type Move struct {
	Color color.Color
	Point *point.Point
}

// BlackMove creates a black move from a point.
func BlackMove(pt *point.Point) *Move {
	return &Move{
		Color: color.Black,
		Point: pt,
	}
}

// BlackMoveList creates a black slice from a point.
func BlackMoveList(pts []*point.Point) []*Move {
	moves := make([]*Move, len(pts))
	for i := range pts {
		moves[i] = BlackMove(pts[i])
	}
	return moves
}

// WhiteMove creates a white move from a point.
func WhiteMove(pt *point.Point) *Move {
	return &Move{
		Color: color.White,
		Point: pt,
	}
}

// WhiteMoveList creates a black slice from a point.
func WhiteMoveList(pts []*point.Point) []*Move {
	moves := make([]*Move, len(pts))
	for i := range pts {
		moves[i] = WhiteMove(pts[i])
	}
	return moves
}
