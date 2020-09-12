// Package point is a basic package for points.
package point

import (
	"bytes"
)

// Point is a basic point. Although simple, the member variables are kept
// private to ensure that Point remains immutable.
type Point struct {
	x int64
	y int64
}

// PointSGFSlce is a slice/Array translation reference between integer
// -Point (index position) and string SGF-Point (byte/char) values
var PointSGFSlce = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h',
	'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u',
	'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
	'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U',
	'V', 'W', 'X', 'Y', 'Z', 'x', 'y', 'z'}

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
// to an SGF Point (two letter string). The returned value is 0-indexed.
func (pt *Point) ToSGF() string {
	sgfOut := ""
	if (pt != nil) && (pt.X() <= 51) && (pt.Y() <= 51) {
		sgfX := string(PointSGFSlce[pt.X()])
		sgfY := string(PointSGFSlce[pt.Y()])
		sgfOut = sgfX + sgfY
	} else {
		panic("Error: *Point entries must not be nil for either" +
			" x y entries and have 0 <= x, y <= 51 values. ")
	}
	return sgfOut
}

// NewFromSGF converts an SGF point (
// two letter string, 0-indexed) to a pointer-type (immutable) *Point.
func NewFromSGF(sgfPt string) *Point {
	var xIndx int
	var yIndx int
	if (sgfPt != "") && (len(sgfPt) == 2) {
		xIndx = bytes.Index(PointSGFSlce, []byte(string(sgfPt[0])))
		yIndx = bytes.Index(PointSGFSlce, []byte(string(sgfPt[1])))
		// x := int64(sgfPt[0]) - aValue
		// y := int64(sgfPt[1]) - aValue
	} else {
		panic("Error: SGF string entries must not be empty and must" +
			" be of length = 2 byte/byte.  ")
		// return New(99, 99)
	}
	return New(int64(xIndx), int64(yIndx))

}