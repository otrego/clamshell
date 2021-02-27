// Package snapshot makes a flattened, abstracted snapshot of a game
// position.
//
// This is based on the flattener logic from glift-core:
// https://github.com/Kashomon/glift-core/tree/master/src/flattener
package snapshot

import (
	"github.com/otrego/clamshell/core/bbox"
	"github.com/otrego/clamshell/core/board"
	"github.com/otrego/clamshell/core/movetree"
)

// Options contains options for creating flattened snapshots.
type Options struct {
	// CropBox allows users to specify a crop-specification.
	CropBox *bbox.CropBox
}

// Create a new Snapshot from a given movetree and path.
func Create(mt *movetree.MoveTree, pos movetree.Path, opts *Options) (*Snapshot, error) {
	size := mt.Root.GameInfo.Size
	n := pos.Apply(mt.Root)
	_, _, err := pos.ApplyToBoard(mt.Root, board.New(size))
	if err != nil {
		return nil, err
	}
	return &Snapshot{
		Comment: n.Comment,
	}, nil
}

// A Snapshot represents a specific board position, which can be used during
// image generation.
type Snapshot struct {
	// Comment for this position.
	Comment string

	// Borad contains the symbolic representation of the display of the board.
	Board *Board
}
