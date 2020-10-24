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

// pointToSgfRef is a translation reference between int64 Point
// and string SGF-Point (rune) values
var pointToSgfRef = map[int64]rune{
	0: 'a', 1: 'b', 2: 'c', 3: 'd', 4: 'e', 5: 'f', 6: 'g',
	7: 'h', 8: 'i', 9: 'j', 10: 'k', 11: 'l', 12: 'm', 13: 'n',
	14: 'o', 15: 'p', 16: 'q', 17: 'r', 18: 's', 19: 't', 20: 'u',
	21: 'v', 22: 'w', 23: 'x', 24: 'y', 25: 'z', 26: 'A', 27: 'B',
	28: 'C', 29: 'D', 30: 'E', 31: 'F', 32: 'G', 33: 'H', 34: 'I',
	35: 'J', 36: 'K', 37: 'L', 38: 'M', 39: 'N', 40: 'O', 41: 'P',
	42: 'Q', 43: 'R', 44: 'S', 45: 'T', 46: 'U', 47: 'V', 48: 'W',
	49: 'X', 50: 'Y', 51: 'Z',
}

// sgfToPointRef is an anonymous function that runs and inverts the
// pointToSgfRef to create the reverse translation reference
var sgfToPointRef = func(in map[int64]rune) map[rune]int64 {
	elem := make(map[rune]int64)
	for key, val := range in {
		elem[val] = key
	}
	return elem
}(pointToSgfRef)

// New creates a new immutable Point.
func New(x, y int64) *Point {
	return &Point{
		x: x,
		y: y,
	}
}

// X returns the x-value
func (pt *Point) X() int64 { return pt.x }

// Y returns the y value
func (pt *Point) Y() int64 { return pt.y }

// ToSGF method receives a pointer argument of type Point,
// and converts it to an SGF Point (two letter string).
func (pt *Point) ToSGF() (string, error) {
	if (pt.X() < 0) || (51 < pt.X()) || (pt.Y() < 0) || 51 < (pt.Y()) {
		return "", fmt.Errorf(
			"*Point int64 x and y value entries must be greater than" +
				" or equal to 0, " +
				"and less than or equal to 51. ")
	}
	return string(pointToSgfRef[pt.X()]) + string(
		pointToSgfRef[pt.Y()]), nil

}

// String method performs to represent and print a Point
func (pt Point) String() string {
	return fmt.Sprintf("{%d,%d}", pt.x, pt.y)
}

// NewFromSGF function receives a pointer argument of type string,
// and converts to an Point reference.
func NewFromSGF(sgfPt string) (*Point, error) {
	if (sgfPt == "") || (len(sgfPt) != 2) {
		return nil, fmt.Errorf(
			"SGF string x and y value entries must be non-empty and" +
				" of length 2 (runes/chars). ")
	}
	return New(sgfToPointRef[rune(sgfPt[0])],
		sgfToPointRef[rune(sgfPt[1])]), nil

}
