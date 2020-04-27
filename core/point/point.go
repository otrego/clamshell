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

// X returns the x-value.
func (p *Point) X() int64 { return p.x }

// Y returns the y-value.
func (p *Point) Y() int64 { return p.y }
