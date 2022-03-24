package snapshot

import (
	"github.com/otrego/clamshell/core/bbox"
	"github.com/otrego/clamshell/core/board"
	"github.com/otrego/clamshell/core/snapshot/symbol"
)

// createBoard creates a Board snapshot from some board state
func createBoard(b *board.Board, cbox *bbox.CropBox) (*Board, error) {
	fb := b.FullBoardState()
	intz := make([][]*Intersection, cbox.BBox.Width())
	for r := cbox.BBox.Top(); r < cbox.BBox.Height(); r++ {
		row := fb[r]
		intz[r] = make([]*Intersection, cbox.BBox.Width())
		for c := cbox.BBox.Left(); c < cbox.BBox.Width(); c++ {
			col := row[c]
			intz[r][c] = &Intersection{
				Stone: symbol.StoneFromColor(col),
			}
		}
	}
	return &Board{
		Intersections: intz,
		CropBox:       cbox,
	}, nil
}

// Board is a snapshot-board representation. It can be cropped.
type Board struct {
	// Intersections contain the intersections for the board.
	Intersections [][]*Intersection

	// The CropBox for the board, specifying the original size and the cropped
	// bounding box.
	CropBox *bbox.CropBox
}
