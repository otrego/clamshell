// Package point is a basic package for points.
package point

// Point is a basic point. Although simple, the member variables are kept
// private to ensure that Point remains immutable.
type Point struct {
	x int64
	y int64
}

// PointToSgfMap is a translation reference between int64 Point
// and string SGF-Point (rune) values

var PointToSgfMap = map[int64]rune{
	0: 'a', 1: 'b', 2: 'c', 3: 'd', 4: 'e', 5: 'f', 6: 'g',
	7: 'h', 8: 'i', 9: 'j', 10: 'k', 11: 'l', 12: 'm', 13: 'n',
	14: 'o', 15: 'p', 16: 'q', 17: 'r', 18: 's', 19: 't', 20: 'u',
	21: 'v', 22: 'w', 23: 'x', 24: 'y', 25: 'z', 26: 'A', 27: 'B',
	28: 'C', 29: 'D', 30: 'E', 31: 'F', 32: 'G', 33: 'H', 34: 'I',
	35: 'J', 36: 'K', 37: 'L', 38: 'M', 39: 'N', 40: 'O', 41: 'P',
	42: 'Q', 43: 'R', 44: 'S', 45: 'T', 46: 'U', 47: 'V', 48: 'W',
	49: 'X', 50: 'Y', 51: 'Z',
}

// SgfToPointMap is a translation reference between string SGF-Point
// (rune) values and int64 Point values
var SgfToPointMap = map[rune]int64{
	'a': 0, 'b': 1, 'c': 2, 'd': 3, 'e': 4, 'f': 5, 'g': 6, 'h': 7,
	'i': 8, 'j': 9, 'k': 10, 'l': 11, 'm': 12, 'n': 13, 'o': 14,
	'p': 15, 'q': 16, 'r': 17, 's': 18, 't': 19, 'u': 20, 'v': 21,
	'w': 22, 'x': 23, 'y': 24, 'z': 25, 'A': 26, 'B': 27, 'C': 28,
	'D': 29, 'E': 30, 'F': 31, 'G': 32, 'H': 33, 'I': 34, 'J': 35,
	'K': 36, 'L': 37, 'M': 38, 'N': 39, 'O': 40, 'P': 41, 'Q': 42,
	'R': 43, 'S': 44, 'T': 45, 'U': 46, 'V': 47, 'W': 48, 'X': 49,
	'Y': 50, 'Z': 51,
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
// The returned value is 0-indexed.
func (pt *Point) ToSGF() string {
	sgfOut := ""
	if (pt.X() <= 51) && (pt.Y() <= 51) {
		sgfX := string(PointToSgfMap[pt.X()])
		sgfY := string(PointToSgfMap[pt.Y()])
		sgfOut = sgfX + sgfY
	} else {
		sgfOut = "--"
	}
	return sgfOut
}

// NewFromSGF converts an SGF point (two letter string)
// to a pointer-type (immutable) *Point.
func NewFromSGF(sgfPt string) *Point {
	var intX int64
	var intY int64
	if (sgfPt != "") && (sgfPt != "--") && (len(sgfPt) == 2) {
		intX = SgfToPointMap[rune(sgfPt[0])]
		intY = SgfToPointMap[rune(sgfPt[1])]
	} else {
		intX = 99
		intY = 99
	}
	return New(intX, intY)

}
