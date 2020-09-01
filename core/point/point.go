// Package point is a basic package for points.
package point

// Point is a basic point. Although simple, the member variables are kept
// private to ensure that Point remains immutable.
type Point struct {
	x int64
	y int64
}

// New creates a new immutable Point.
func New(x, y int64) *Point {
	return &Point{
		x: x,
		y: y,
	}
}

type SGF struct {
	x string
	y string
}

// X returns the x-value.
func (p *Point) X() int64 { return p.x }

// Y returns the y-value.
func (p *Point) Y() int64 { return p.y }

// Two Go-game board position/point translation methods
// 1. '*Point' immutable int point, to SGF(String of two letters) point
func (pt *Point) ToSGF() string {
	sgfX := string(rune((pt.X()) + 97))
	sgfY := string(rune((pt.Y()) + 97))
	// two lower case letters
	return sgfX + sgfY
}

//[... worked to here]

// 3. 'p' SGF(two letter string) Point, to '*Point' immutable Int Point
/*func translateToINTPt01(sgfPt string) *Point {
	sgfPtX := sgfPt[0]
	sgfPtY := sgfPt[1]
	x := int(sgfPtX) - 97
	y := int(sgfPtY) - 97
	return New(int64(x), int64(y))
}*/
