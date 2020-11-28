package problems

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/otrego/clamshell/core/board"
	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/sgf"
)

// Flatten
func Flatten(tp movetree.Treepath, g *movetree.MoveTree) (*movetree.MoveTree, error) {
	b := populateBoard(tp, g)
	placements := b.GetPlacements()

	var sbSGF strings.Builder
	var sbAB strings.Builder
	var sbAW strings.Builder

	sbAB.WriteString("AB")
	sbAW.WriteString("AW")
	for _, move := range placements {
		str, err := move.Point().ToSGF()
		if err != nil {
			log.Fatal(err)
		}
		if move.Color() == color.Black {
			sbAB.WriteString(fmt.Sprintf("[%s]", str))
		} else {
			sbAW.WriteString(fmt.Sprintf("[%s]", str))
		}
	}

	sbSGF.WriteString("(;")
	for key, value := range g.Root.Properties {
		if key != "AW" && key != "AB" {
			sbSGF.WriteString(fmt.Sprintf("%s[%s]", key, value))
		}
	}
	sbSGF.WriteString(sbAB.String())
	sbSGF.WriteString(sbAW.String())
	sbSGF.WriteString(")")

	return sgf.Parse(sbSGF.String())

}

func populateBoard(tp movetree.Treepath, g *movetree.MoveTree) *board.Board {
	n := g.Root
	i, err := strconv.Atoi(n.Properties["SZ"][0])
	if err != nil {

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
			log.Fatal(fmt.Errorf("treepath leads to nil movetree node"))
		}

		_, err := b.PlaceStone(n.Move)
		if err != nil {
			log.Fatal(err)
		}
	}
	return b
}
