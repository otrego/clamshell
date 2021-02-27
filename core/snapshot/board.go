package snapshot

import (
	"github.com/otrego/clamshell/core/bbox"
	"github.com/otrego/clamshell/core/board"
	"github.com/otrego/clamshell/core/snapshot/symbol"
)

// createBoard creates a Board snapshot from some board state
func createBoard(b *board.Board, cbox *bbox.CropBox) (*Board, error) {
	fb := b.FullBoardState()
	intz := make([][]*Intersection, len(fb))
	for i, row := range fb {
		intz[i] = make([]*Intersection, len(row[0]))
		for j, col := range row {
			intz[i][j] = &Intersection{
				Stone: symbol.StoneFromColor(col),
			}
		}
	}
	return &Board{
		Intersections: intz,
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
