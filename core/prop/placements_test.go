package prop

import (
	"testing"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/move"
	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/point"
)

func TestConvertFromSGF_Placements(t *testing.T) {
	testCases := []fromSGFTestCase{
		{
			desc: "black placements",
			prop: "AB",
			data: []string{"aa", "bb"},
			makeExpNode: func(n *movetree.Node) {
				n.Placements = []*move.Move{
					move.New(color.Black, point.New(0, 0)),
					move.New(color.Black, point.New(1, 1)),
				}
			},
		},
		{
			desc: "white placements",
			prop: "AW",
			data: []string{"aa", "bb"},
			makeExpNode: func(n *movetree.Node) {
				n.Placements = []*move.Move{
					move.New(color.White, point.New(0, 0)),
					move.New(color.White, point.New(1, 1)),
				}
			},
		},
	}

	testConvertFromSGFCases(t, testCases)
}

func TestConvertNode_Placements(t *testing.T) {
	testCases := []convertNodeTestCase{
		{
			desc: "black placements",
			makeNode: func(n *movetree.Node) {
				n.Placements = []*move.Move{
					move.New(color.Black, point.New(0, 1)),
					move.New(color.Black, point.New(0, 2)),
				}
			},
			expOut: "AB[ab][ac]",
		},
		{
			desc: "white placements",
			makeNode: func(n *movetree.Node) {
				n.Placements = []*move.Move{
					move.New(color.White, point.New(0, 1)),
					move.New(color.White, point.New(0, 2)),
				}
			},
			expOut: "AW[ab][ac]",
		},
	}

	testConvertNodeCases(t, testCases)
}
