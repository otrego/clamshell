// Package point is a basic package for points.
package point

import (
	"fmt"
)

// Point is a basic point. Although simple, the member variables are kept
// private to ensure that Point remains immutable.
type Point struct {
	x int
	y int
}

// New creates a new immutable Point.
func New(x, y int) *Point {
	return &Point{
		x: x,
		y: y,
	}
}

// X returns the x-value.
func (pt *Point) X() int { return pt.x }

// Y returns the y-value.
func (pt *Point) Y() int { return pt.y }

// ToSGF converts a pointer-type (immutable) *Point
// to an SGF Point (two letter string).
func (pt *Point) ToSGF() (string, error) {
	return toSGF(pt)
}

// Equal returns whether this point is equal to another point.
func (pt *Point) Equal(other *Point) bool {
	return other != nil && pt.X() == other.X() && pt.Y() == other.Y()
}

// String converts to string representation of a Point.
func (pt *Point) String() string {
	return fmt.Sprintf("{%d,%d}", pt.x, pt.y)
}

// Key is a convenience helper to convert this point to a key-struct.
func (pt *Point) Key() Key {
	return Key{X: pt.X(), Y: pt.Y()}
}

// Key is a point-struct that is used for keys in maps. As such, it's intended
// to be used like the following:
//
//     Key{X:12, Y:15}
type Key struct {
	X int
	Y int
}

// Point converts a point-Key back to a point.
func (k Key) Point() *Point {
	return New(k.X, k.Y)
}
