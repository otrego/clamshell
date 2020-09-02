// Package point is a basic package for points.
package point

// Point is a basic point. Although simple, the member variables are kept
// private to ensure that Point remains immutable.
type Point struct {
	x int64
	y int64
}

const aValue = int64('a')

// New creates a new immutable Point.
func New(x, y int64) *Point {
	return &Point{
		x: x,
		y: y,
	}
}

// X returns the x-value.
func (p *Point) X() int64 { return p.x }

// Y returns the y-value.
func (p *Point) Y() int64 { return p.y }

// Two Go-game board position/point translation methods
// 1. '*Point' immutable int point, to SGF(String of two letters) point
func (pt *Point) ToSGF() string {
	sgfX := string(rune((pt.X()) + aValue))
	sgfY := string(rune((pt.Y()) + aValue))
	// two lower case letters
	return sgfX + sgfY
}

// 3. 'p' SGF(two letter string) Point, to '*Point' immutable Int Point
func NewFromSGF(sgfPt string) *Point {
	sgfX := sgfPt[0]
	sgfY := sgfPt[1]
	x := int64(sgfX) - aValue
	y := int64(sgfY) - aValue
	return New(int64(x), int64(y))
}
