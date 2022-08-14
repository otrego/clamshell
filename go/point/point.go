// Package point is a basic package for points.
package point

import (
	"encoding/json"
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

// pointInternal is an internal struct for the purposes of JSON Conversion.
type pointInternal struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// MarshalJSON indicates to the JSON library how to marshal a Point.
func (pt *Point) MarshalJSON() ([]byte, error) {
	return json.Marshal(&pointInternal{
		X: pt.x,
		Y: pt.y,
	})
}

// UnmarshalJSON indicates to the JSON library how to unmarshal a Point.
func (pt *Point) UnmarshalJSON(data []byte) error {
	var pti pointInternal
	if err := json.Unmarshal(data, &pti); err != nil {
		return err
	}
	pt.x = pti.X
	pt.y = pti.Y
	return nil
}
