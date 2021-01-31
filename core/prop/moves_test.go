package prop

import (
	"testing"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/move"
	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/point"
)

func TestConvertFromSGF_Moves(t *testing.T) {
	testCases := []fromSGFTestCase{
		{
			desc: "black move: pass",
			prop: "B",
			data: []string{},
			makeExpNode: func(n *movetree.Node) {
				n.Move = move.NewPass(color.Black)
			},
		},
		{
			desc: "black move: pass empty str",
			prop: "B",
			data: []string{""},
			makeExpNode: func(n *movetree.Node) {
				n.Move = move.NewPass(color.Black)
			},
		},
		{
			desc: "black move",
			prop: "B",
			data: []string{"ab"},
			makeExpNode: func(n *movetree.Node) {
				n.Move = move.New(color.Black, point.New(0, 1))
			},
		},
		{
			desc: "white move",
			prop: "W",
			data: []string{"ab"},
			makeExpNode: func(n *movetree.Node) {
				n.Move = move.New(color.White, point.New(0, 1))
			},
		},
	}

	testConvertFromSGFCases(t, testCases)
}

func TestConvertNode_Moves(t *testing.T) {
	testCases := []convertNodeTestCase{
		{
			desc: "black move: pass",
			makeNode: func(n *movetree.Node) {
				n.Move = move.NewPass(color.Black)
			},
			expOut: "B[]",
		},
		{
			desc: "black move: non-pass",
			makeNode: func(n *movetree.Node) {
				n.Move = move.New(color.Black, point.New(0, 1))
			},
			expOut: "B[ab]",
		},
		{
			desc: "white move: non-pass",
			makeNode: func(n *movetree.Node) {
				n.Move = move.New(color.White, point.New(0, 1))
			},
			expOut: "W[ab]",
		},
	}

	testConvertNodeCases(t, testCases)
}
