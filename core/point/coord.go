// Package point is a basic package for points.
package point

import (
	"fmt"
)

// Coord is a float point for representing points in space
type Coord struct {
	x float64
	y float64
}

// NewCoord creates a new immutable Coord.
func NewCoord(x, y float64) *Coord {
	return &Coord{
		x: x,
		y: y,
	}
}

// X returns the x-value.
func (pt *Coord) X() float64 { return pt.x }

// Y returns the y-value.
func (pt *Coord) Y() float64 { return pt.y }

// String method to represent and print a Coord.
func (pt *Coord) String() string {
	return fmt.Sprintf("{%v,%v}", pt.x, pt.y)
}

// Equal returns whether this point is equal to another point.
func (pt *Coord) Equal(other *Coord) bool {
	return other != nil && pt.X() == other.X() && pt.Y() == other.Y()
}
