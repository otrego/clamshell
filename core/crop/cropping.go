package crop

import (
	"github.com/otrego/clamshell/core/point"
)

//CroppingPreset : Convenience enum for specifying a cropping direction.
type CroppingPreset int64

const (
	//TopLeft see below
	// X -
	// - -
	TopLeft CroppingPreset = iota

	//TopRight see below
	// - X
	// - -
	TopRight

	//BottomRight see below
	// - -
	// - X
	BottomRight

	//BottomLeft see below
	// - -
	// X -
	BottomLeft

	//Top see below
	// X X
	// - -
	Top

	//Left see below
	// X -
	// X -
	Left

	//Right see below
	// - X
	// - X
	Right

	//Bottom see below
	// - -
	// X X
	Bottom

	//All see below
	// X X
	// X X
	All
)

//Cropping is a bounding box, specified by intersection points.
type Cropping struct {
	TopLeft  *point.Point
	BotRight *point.Point
}

//FromPreset : Create a cropping box from the maxInts. Note that the integer points in the
// crop box are 0 indexed, but maxInts is 1-indexed.  In other words, we would
// typically expect the max ints to range from 9 to 19.
//
// Following the SGF covention, we consider the topleft to be 0,0
func FromPreset(p CroppingPreset, maxInts int64) *Cropping {
	halfInts := maxInts / 2
	var minInts int64 = 0
	top := minInts
	left := minInts
	bot := maxInts
	right := maxInts
	switch p {
	case All: // nothing to change
	case Left:
		right = halfInts + 1
	case Right:
		left = halfInts - 1
	case Top:
		bot = halfInts + 1
	case Bottom:
		top = halfInts - 1
	case TopLeft:
		bot = halfInts + 1
		right = halfInts + 2
	case TopRight:
		bot = halfInts + 1
		left = halfInts - 2
	case BottomLeft:
		top = halfInts - 1
		right = halfInts + 2
	case BottomRight:
		top = halfInts - 1
		left = halfInts - 2
	}
	return &Cropping{
		point.New(top, left),
		point.New(bot, right),
	}
}
