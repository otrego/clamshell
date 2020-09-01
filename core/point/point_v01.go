// Package point is a basic package for points.
package point

/*import (
	"strings"
)

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

// X returns the x-value.
func (p *Point) X() int64 { return p.x }

// Y returns the y-value.
func (p *Point) Y() int64 { return p.y }

// Kilo Foxtrot addition below
// Three point translation methods, one is an alternate

// 1. '*Point' immutable int point, to SGF(String of two letters) Point
func ToSGF(pt *Point) string {
	sgfPtX := string(rune((pt.X()) + 97))
	sgfPtY := string(rune((pt.Y()) + 97))
	// two lower case letters
	// return sgfPtX + sgfPtY
	return strings.ToLower(sgfPtX) + strings.ToLower(sgfPtY)
}

// Alternate, in case the incoming point is not immutable
// 2. 'Point' mutable int point, to SGF(String of two letters) Point
func translateToSGFPt02(pt Point) string {
	sgfPtX := string(rune((pt.X()) + 97))
	sgfPtY := string(rune((pt.Y()) + 97))
	// two lower case letters
	// return sgfPtX + sgfPtY
	return strings.ToLower(sgfPtX) + strings.ToLower(sgfPtY)
}

// 3. 'p' SGF(two letter string) Point, to '*Point' immutable Int Point
func translateToINTPt01(sgfPt string) *Point {
	sgfPtX := sgfPt[0]
	sgfPtY := sgfPt[1]
	x := int(sgfPtX) - 97
	y := int(sgfPtY) - 97
	return New(int64(x), int64(y))
}*/
