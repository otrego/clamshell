// Package bbox provides utilities for managing various types of bounding boxes.
package bbox

import (
	"fmt"

	"github.com/otrego/clamshell/core/point"
)

// New creates a new bounding box.
func New(tl, br *point.Point) (*BoundingBox, error) {
	if tl.X() >= br.X() || tl.Y() >= br.Y() {
		// Open question: Should we allow 0-width / 0-height bboxes? It doesn't seem
		// to have much practical value.
		return nil, fmt.Errorf("bottom right must be greater than topleft "+
			"but top left was %v and bottom right was %v", tl, br)
	}
	return &BoundingBox{
		tl: tl,
		br: br,
	}, nil
}

// BoundingBox is a bounding box for a board, which provides a top-left and
// bottom-right.
type BoundingBox struct {
	tl *point.Point
	br *point.Point
}

// TopLeft point of the bounding box.
func (b *BoundingBox) TopLeft() *point.Point { return b.tl }

// BotRight point of the bounding box.
func (b *BoundingBox) BotRight() *point.Point { return b.br }

// Top side of the bounding box.
func (b *BoundingBox) Top() int64 { return b.tl.Y() }

// Left side of the bounding box
func (b *BoundingBox) Left() int64 { return b.tl.X() }

// Bottom side of the bounding box.
func (b *BoundingBox) Bottom() int64 { return b.br.Y() }

// Right side of the bounding box.
func (b *BoundingBox) Right() int64 { return b.br.X() }

// Width of the bounding box.
func (b *BoundingBox) Width() int64 { return b.br.X() - b.tl.X() }

// Height of the bounding box.
func (b *BoundingBox) Height() int64 { return b.br.Y() - b.tl.Y() }
