package bbox

import (
	"fmt"

	"github.com/otrego/clamshell/core/point"
)

// A CropBox is a bounding box that contains a strict subset of a board.
type CropBox struct {
	BBox         *BoundingBox
	OriginalSize int
}

// CroppingPreset is a convenience enum for specifying a cropping direction.
type CroppingPreset int

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

// CropBoxFromPreset creates a cropping box from the original board size
// (typically 9, 13, 19). Note that the integer points in the crop box are 0
// indexed, but originalSize is 1-indexed.  In other words, we would typically
// expect the max ints to range from 9 to 19.
//
// Following the SGF covention, we consider the topleft to be 0,0
func CropBoxFromPreset(p CroppingPreset, boardSize int) (*CropBox, error) {
	bs := boardSize

	halfInts := bs / 2
	var minInts int

	top := minInts
	left := minInts
	bot := bs
	right := bs

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
	bb, err := New(point.New(top, left), point.New(bot, right))
	if err != nil {
		return nil, fmt.Errorf("error while creating cropbox: %v", err)
	}
	return &CropBox{
		BBox:         bb,
		OriginalSize: boardSize,
	}, nil
}
