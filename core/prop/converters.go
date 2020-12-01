package prop

import (
	"fmt"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/move"
	"github.com/otrego/clamshell/core/movetree"
)

// converters contain all the property converters.
var converters = []*Converter{
	// Black Placements
	&Converter{
		Prop: "AB",
		From: func(n *movetree.Node, data []string) error {
			moves, err := move.ListFromSGFPoints(color.Black, data)
			if err != nil {
				return err
			}
			n.Placements = append(n.Placements, moves...)
			return nil
		},
		To: func(n *movetree.Node) ([]string, error) {
			var out []string
			for _, mv := range n.Placements {
				if mv.Color() != color.Black {
					continue
				}
				sgfPt, err := mv.Point().ToSGF()
				if err != nil {
					return nil, err
				}
				out = append(out, sgfPt)
			}
			return out, nil
		},
	},

	// White Placements
	&Converter{
		Prop: "AW",
		From: func(n *movetree.Node, data []string) error {
			moves, err := move.ListFromSGFPoints(color.White, data)
			if err != nil {
				return err
			}
			n.Placements = append(n.Placements, moves...)
			return nil
		},
		To: func(n *movetree.Node) ([]string, error) {
			var out []string
			for _, mv := range n.Placements {
				if mv.Color() != color.White {
					continue
				}
				sgfPt, err := mv.Point().ToSGF()
				if err != nil {
					return nil, err
				}
				out = append(out, sgfPt)
			}
			return out, nil
		},
	},

	// Black Moves
	&Converter{
		Prop: "B",
		From: func(n *movetree.Node, data []string) error {
			if n.Move != nil {
				return fmt.Errorf("found two moves on one node")
			}
			if len(data) != 1 && len(data) != 0 {
				return fmt.Errorf("expected black move data to have exactly one value or zero values")
			}
			if len(data) == 0 {
				data = []string{""}
			}
			move, err := move.FromSGFPoint(color.Black, data[0])
			if err != nil {
				return err
			}
			n.Move = move
			return nil
		},
		To: func(n *movetree.Node) ([]string, error) {
			mv := n.Move
			if mv.Color() != color.Black {
				return nil, nil
			}
			if mv.IsPass() {
				// Return non-nil slice to indicate it should be stored.
				return []string{}, nil
			}
			sgfPt, err := mv.Point().ToSGF()
			if err != nil {
				return nil, err
			}
			return []string{sgfPt}, nil
		},
	},

	// White Moves
	&Converter{
		Prop: "W",
		From: func(n *movetree.Node, data []string) error {
			if n.Move != nil {
				return fmt.Errorf("found two moves on one node")
			}
			if len(data) != 1 && len(data) != 0 {
				return fmt.Errorf("expected white move data to have exactly one or zero values")
			}
			if len(data) == 0 {
				data = []string{""}
			}
			move, err := move.FromSGFPoint(color.Black, data[0])
			if err != nil {
				return err
			}
			n.Move = move
			return nil
		},
		To: func(n *movetree.Node) ([]string, error) {
			mv := n.Move
			if mv.Color() != color.Black {
				return nil, nil
			}
			if mv.IsPass() {
				// Return non-nil slice to indicate it should be stored as such.
				return []string{}, nil
			}
			sgfPt, err := mv.Point().ToSGF()
			if err != nil {
				return nil, err
			}
			return []string{sgfPt}, nil
		},
	},
}

var convMap = func(conv []*Converter) map[Prop]*Converter {
	mp := make(map[Prop]*Converter)
	for _, c := range conv {
		mp[c.Prop] = c
	}
	return mp
}(converters)

// HasConverter indicates whether there's a known SGF Property converter.
func HasConverter(prop string) bool {
	_, ok := convMap[Prop(prop)]
	return ok
}

// PropConverter gets a property converter, returning nil if no property
// converter can be found.
func PropConverter(prop string) *Converter {
	return convMap[Prop(prop)]
}
