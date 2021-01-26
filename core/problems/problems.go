package problems

import (
	"fmt"

	"github.com/otrego/clamshell/core/board"
	"github.com/otrego/clamshell/core/movetree"
)

// Flatten takes a Pathand a MoveTree and returns
// the root of the flat movetree. the flat movetree
// ignores the last two nodes of the treepath
func Flatten(tp movetree.Path, g *movetree.MoveTree) (*movetree.MoveTree, error) {
	b, err := PopulateBoard(tp, g)
	if err != nil {
		return nil, err
	}

	gflat := movetree.New()
	gflat.Root.Placements = b.GetFullBoardState()
	gflat.Root.GameInfo.Size = g.Root.GameInfo.Size

	for key, value := range g.Root.SGFProperties {
		gflat.Root.SGFProperties[key] = value
	}

	return gflat, nil
}

// PopulateBoard populates a go board given a MoveTree and Path
// Captures are intentionally discarded here.
// returns the populated board.
func PopulateBoard(tp movetree.Path, g *movetree.MoveTree) (*board.Board, error) {
	n := g.Root
	i := n.GameInfo.Size

	b := board.NewBoard(i)
	for _, move := range n.Placements {
		b.PlaceStone(move)
	}

	// tp ends one move after the blunder, so we follow the treepath
	// to the move right before the blunder using len(tp) - 2;
	for i := 0; i < len(tp)-2; i++ {
		n = n.Next(tp[i])
		if n == nil {
			return nil, fmt.Errorf("treepath leads to nil movetree node")
		}

		_, err := b.PlaceStone(n.Move)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}
