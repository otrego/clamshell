// Package snapshot makes a flattened, abstracted snapshot of a game
// position.
//
// This is based on the flattener logic from glift-core:
// https://github.com/Kashomon/glift-core/tree/master/src/flattener
package snapshot

import (
	"github.com/otrego/clamshell/core/bbox"
	"github.com/otrego/clamshell/core/movetree"
)

// Options contains options for creating flattened snapshots.
type Options struct {
}

// Create a new Snapshot from a given movetree and path.
func Create(mt *movetree.MoveTree, position movetree.Path, opts *Options) *Snapshot {
	return &Snapshot{}
}

// A Snapshot represents a specific board position, which can be used during
// image generation.
type Snapshot struct {
	// Comment for this position.
	Comment string

	// Borad contains the symbolic representation of the display of the board.
	Board *Board
}

// Board is a snapshot-board representation. It can be cropped.
type Board struct {
	// Intersections contain the intersections for the board.
	Intersections [][]*Intersection

	// The CropBox for the board, specifying the original size and the cropped
	// bounding box.
	CropBox *bbox.CropBox
}
