package problems

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/otrego/clamshell/core/board"
	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/sgf"
)

// Flatten takes a Treepath and a MoveTree and returns
// the root of the flat movetree. the flat movetree
// ignores the last two nodes of the treepath
func Flatten(tp movetree.Treepath, g *movetree.MoveTree) (*movetree.MoveTree, error) {
	b, err := PopulateBoard(tp, g)
	placements := b.GetPlacements()

	if err != nil {
		return nil, err
	}

	var sbSGF strings.Builder
	var sbAB strings.Builder
	var sbAW strings.Builder

	// write new SGF AB and AW properties and values
	// for the new Flattened SGF file
	sbAB.WriteString("AB")
	sbAW.WriteString("AW")
	for _, move := range placements {
		str, err := move.Point().ToSGF()
		if err != nil {
			return nil, err
		}
		if move.Color() == color.Black {
			sbAB.WriteString(fmt.Sprintf("[%s]", str))
		} else {
			sbAW.WriteString(fmt.Sprintf("[%s]", str))
		}
	}

	// build new SGF file and replace the old AW and AB
	// with the new AW and BW
	sbSGF.WriteString("(;")
	for key, value := range g.Root.Properties {
		if key != "AW" && key != "AB" {
			sbSGF.WriteString(fmt.Sprintf("%s%s", key, value))
		}
	}
	sbSGF.WriteString(sbAB.String())
	sbSGF.WriteString(sbAW.String())
	sbSGF.WriteString(")")

	// create new Flattened movetree
	gFlat, err := sgf.Parse(sbSGF.String())
	if err != nil {
		return nil, err
	}

	return gFlat, nil

}

// PopulateBoard populates a go board given a MoveTree and Treepath
// returns the populated board.
func PopulateBoard(tp movetree.Treepath, g *movetree.MoveTree) (*board.Board, error) {
	n := g.Root
	i, err := strconv.Atoi(n.Properties["SZ"][0])
	if err != nil {
		return nil, err
	}
	b := board.NewBoard(i)
	for _, move := range n.Placements {
		b.AddStone(move)
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
