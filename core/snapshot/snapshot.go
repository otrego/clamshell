// Package snapshot makes a flattened, abstracted snapshot of a particular game
// position.
//
// This is based on the flattener logic from glift-core:
// https://github.com/Kashomon/glift-core/tree/master/src/flattener
package snapshot

import "github.com/otrego/clamshell/core/movetree"

// Options contains options for creating snapshots.
type Options struct {
}

// New creates a new Snapshot from a given movetree and path.
func New(mt *movetree.MoveTree, position movetree.Path, opts *Options) *Snapshot {
	return &Snapshot{
		OriginalSize: mt.Root.GameInfo.Size,
	}
}

// A snapshot represents a specific board position, which can be used during
// image generation.
//
// **Snapshots are intended to be immutable.
type Snapshot struct {
	// OriginalSize of the Go-Board. The actual size is likely less due to
	// cropping.
	OriginalSize int
}

// Board is a snapshot-board representation. It is different than the game-logic
// board in that it is frequently cropped and designed for rendering.
type Board struct {
	Intersections [][]*Intersection
}

type Intersection struct {
}
