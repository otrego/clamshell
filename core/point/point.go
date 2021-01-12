// Package point is a basic package for points.
package point

import (
	"fmt"
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
func (pt *Point) X() int64 { return pt.x }

// Y returns the y-value.
func (pt *Point) Y() int64 { return pt.y }

// ToSGF converts a pointer-type (immutable) *Point
// to an SGF Point (two letter string).
func (pt *Point) ToSGF() (string, error) {
	return toSGF(pt)
}

// Equals returns whether this point is equal to another point.
func (pt *Point) Equals(other *Point) bool {
	return pt.X() == other.X() && pt.Y() == other.Y()
}

// String converts to string representation of a Point.
func (pt *Point) String() string {
	return fmt.Sprintf("{%d,%d}", pt.x, pt.y)
}

// Convenience helper to convert this point to a key-struct.
func (pt *Point) Key() Key {
	return Key{X: pt.X(), Y: pt.Y()}
}

// Key is a point-struct that is used for keys in maps. As such, it's intended
// to be used like the following:
//
//     Key{X:12, Y:15}
type Key struct {
	X int64
	Y int64
}

// ToPoint converts a point-Key back to a point.
func (k Key) Point() *Point {
	return New(k.X, k.Y)
}
