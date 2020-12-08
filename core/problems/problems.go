package problems

import (
	"fmt"
	"strconv"

	"github.com/otrego/clamshell/core/board"
	"github.com/otrego/clamshell/core/movetree"
)

// Flatten takes a Treepath and a MoveTree and returns
// the root of the flat movetree. the flat movetree
// ignores the last two nodes of the treepath
func Flatten(tp movetree.Treepath, g *movetree.MoveTree) (*movetree.MoveTree, error) {
	b, err := PopulateBoard(tp, g)
	if err != nil {
		return nil, err
	}

	gflat := movetree.New()
	gflat.Root.Placements = b.GetFullBoardState()

	for key, value := range g.Root.Properties {
		gflat.Root.Properties[key] = value
	}

	return gflat, nil

}

// PopulateBoard populates a go board given a MoveTree and Treepath
// Captures are intentionally discarded here.
// returns the populated board.
func PopulateBoard(tp movetree.Treepath, g *movetree.MoveTree) (*board.Board, error) {
	n := g.Root
	i, err := strconv.Atoi(n.Properties["SZ"][0])
	if err != nil {
		return nil, err
	}
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

		_, err2 := b.PlaceStone(n.Move)
		if err2 != nil {
			return nil, err
		}
	}
	return b, nil
}
