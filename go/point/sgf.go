package point

import (
	"errors"
	"fmt"
)

var SGFConversionErr = errors.New("error converting point from sgf pointF")

// pointToSgfMap is a translation reference between int Point
// and string SGF-Point (rune) values
var pointToSgfMap = map[int]rune{
	0: 'a', 1: 'b', 2: 'c', 3: 'd', 4: 'e', 5: 'f', 6: 'g',
	7: 'h', 8: 'i', 9: 'j', 10: 'k', 11: 'l', 12: 'm', 13: 'n',
	14: 'o', 15: 'p', 16: 'q', 17: 'r', 18: 's', 19: 't', 20: 'u',
	21: 'v', 22: 'w', 23: 'x', 24: 'y', 25: 'z', 26: 'A', 27: 'B',
	28: 'C', 29: 'D', 30: 'E', 31: 'F', 32: 'G', 33: 'H', 34: 'I',
	35: 'J', 36: 'K', 37: 'L', 38: 'M', 39: 'N', 40: 'O', 41: 'P',
	42: 'Q', 43: 'R', 44: 'S', 45: 'T', 46: 'U', 47: 'V', 48: 'W',
	49: 'X', 50: 'Y', 51: 'Z',
}

// sgfToPointMap is a translation reference between string SGF-Point
// (rune) values and int Point values
var sgfToPointMap = func(m map[int]rune) map[rune]int {
	out := make(map[rune]int)
	for key, val := range m {
		out[val] = key
	}
	return out
}(pointToSgfMap)

// NewFromSGF converts an SGF point (two letter string) to a Point.
func NewFromSGF(sgfPt string) (*Point, error) {
	if sgfPt == "" || len(sgfPt) != 2 {
		return nil, fmt.Errorf("%w sgf point must be non-empty and two letter char, but was %s", SGFConversionErr, sgfPt)
	}
	intX, ok := sgfToPointMap[rune(sgfPt[0])]
	if !ok {
		return nil, fmt.Errorf("%w could not convert coordinate for x-value of sgf point %s; only a-zA-Z (minus i/I) are allowed", SGFConversionErr, sgfPt)
	}
	intY, ok := sgfToPointMap[rune(sgfPt[1])]
	if !ok {
		return nil, fmt.Errorf("%w could not convert coordinate for y-value of sgf point %s; only a-zA-Z (minus i/I) are allowed", SGFConversionErr, sgfPt)
	}
	return New(intX, intY), nil
}

// toSGF converts a point to an SGF coordinate.
func toSGF(pt *Point) (string, error) {
	if pt.X() < 0 || pt.X() > 51 || pt.Y() < 0 || pt.Y() > 51 {
		return "", fmt.Errorf("%w error converting point to SGF-point; points must be between 0 and 51 inclusive. found %v", SGFConversionErr, pt.String())
	}
	return string(pointToSgfMap[pt.X()]) + string(pointToSgfMap[pt.Y()]), nil
}
