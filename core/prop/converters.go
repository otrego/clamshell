package prop

import (
	"fmt"
	"sort"
	"strings"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/move"
	"github.com/otrego/clamshell/core/movetree"
)

// Scope indicates the scope for a property.
type Scope string

const (
	// RootScope indicates a property that only applies to the root node.
	RootScope Scope = "RootScope"

	// AllScope indicates a property that applies to all nodes.
	AllScope Scope = "AllScope"
)

// FromSGF converts an SGF Property to node property
type FromSGF func(node *movetree.Node, prop string, values []string) error

// ToSGF converts an Node property to an SGF property list.
type ToSGF func(node *movetree.Node) (string, error)

// A SGFConverter converts SGF properties to / from node properties.
type SGFConverter struct {
	// Props is the name of the SGF properties that apply to this
	// property-converter.
	// Ex: {"AW", "AB"}
	Props []Prop
	// Scope indicates what the property-converter applies to (Root, All)
	Scope Scope
	// From converts from SGF data
	From FromSGF
	// To converts to SGF data
	To ToSGF
}

// converters contain all the property converters.
var converters = []*SGFConverter{
	// Placements
	&SGFConverter{
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
	},

	// Moves
	&SGFConverter{
		Props: []Prop{"B", "W"},
		Scope: AllScope,
		From: func(n *movetree.Node, prop string, data []string) error {
			col, err := color.FromSGFProp(prop)
			if err != nil {
				return err
			}
			if n.Move != nil {
				return fmt.Errorf("found two moves on one node at move")
			}
			if len(data) != 1 && len(data) != 0 {
				return fmt.Errorf("expected black move data to have exactly one value or zero values")
			}
			if len(data) == 0 {
				data = []string{""}
			}
			move, err := move.FromSGFPoint(col, data[0])
			if err != nil {
				return err
			}
			n.Move = move
			return nil
		},
		To: func(n *movetree.Node) (string, error) {
			mv := n.Move
			if mv == nil {
				return "", nil
			}
			var col string
			if mv.Color() == color.Black {
				col = "B"
			} else if mv.Color() == color.White {
				col = "W"
			}
			if mv.IsPass() {
				// Return non-nil slice to indicate it should be stored.
				return col + "[]", nil
			}
			sgfPt, err := mv.Point().ToSGF()
			if err != nil {
				return "", err
			}
			return col + "[" + sgfPt + "]", nil
		},
	},
}

var propToConv = func(conv []*SGFConverter) map[Prop]*SGFConverter {
	mp := make(map[Prop]*SGFConverter)
	for _, c := range conv {
		for _, p := range c.Props {
			mp[p] = c
		}
	}
	return mp
}(converters)

// HasConverter indicates whether there's a known SGF Property converter.
func HasConverter(prop string) bool {
	_, ok := propToConv[Prop(prop)]
	return ok
}

// Converter gets a property converter for converting to/from SGf, returning nil
// if no property converter can be found.
func Converter(prop string) *SGFConverter {
	return propToConv[Prop(prop)]
}

// ConvertNode converts all the properties in a node
func ConvertNode(n *movetree.Node) (string, error) {
	var sb strings.Builder
	for _, c := range converters {
		if c.Scope == RootScope && n.MoveNum() != 0 {
			// skip non-root-scoped properties for non-root nodes.
			continue
		}
		s, err := c.To(n)
		if err != nil {
			return "", err
		}
		sb.WriteString(s)
	}

	var keys []string
	for key := range n.Properties {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		sb.WriteString(key)
		for _, value := range n.Properties[key] {
			sb.WriteString("[" + value + "]")
		}
	}
	return sb.String(), nil
}
