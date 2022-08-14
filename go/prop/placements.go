package prop

import (
	"strings"

	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/move"
	"github.com/otrego/clamshell/go/movetree"
)

// placementsConv converts stone-placements AW, AB.
var placementsConv = &SGFConverter{
	Props: []Prop{"AB", "AW"},
	Scope: AllScope,
	From: func(n *movetree.Node, prop string, data []string) error {
		col, err := color.FromSGFProp(prop)
		if err != nil {
			return err
		}
		moves, err := move.ListFromSGFPoints(col, data)
		if err != nil {
			return err
		}
		n.Placements = append(n.Placements, moves...)
		return nil
	},
	To: func(n *movetree.Node) (string, error) {
		if len(n.Placements) == 0 {
			return "", nil
		}
		var black []string
		var white []string
		for _, mv := range n.Placements {
			sgfPt, err := mv.Point().ToSGF()
			if err != nil {
				return "", err
			}
			if mv.Color() == color.Black {
				black = append(black, sgfPt)
			} else if mv.Color() == color.White {
				white = append(white, sgfPt)
			}
		}
		var out strings.Builder
		if len(black) > 0 {
			out.WriteString("AB")
		}
		for _, pt := range black {
			out.WriteString("[" + pt + "]")
		}
		if len(white) > 0 {
			out.WriteString("AW")
		}
		for _, pt := range white {
			out.WriteString("[" + pt + "]")
		}
		return out.String(), nil
	},
}
