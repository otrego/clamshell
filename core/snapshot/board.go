package snapshot

import (
	"github.com/otrego/clamshell/core/bbox"
	"github.com/otrego/clamshell/core/board"
)

// createBoard creates a Board snapshot from some board state
func createBoard(b *board.Board, cbox *bbox.CropBox) (*Board, error) {
	// TODO(kashomon): Add support for this
	return nil, nil
}

// Board is a snapshot-board representation. It can be cropped.
type Board struct {
	// Intersections contain the intersections for the board.
	Intersections [][]*Intersection

	// The CropBox for the board, specifying the original size and the cropped
	// bounding box.
	CropBox *bbox.CropBox
}
