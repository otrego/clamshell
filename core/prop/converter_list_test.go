package prop

import (
	"reflect"
	"testing"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/errcheck"
	"github.com/otrego/clamshell/core/move"
	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/point"
)

type propGetter func(*movetree.Node) interface{}

type setprops func(*movetree.Node) *movetree.Node

func TestConverters_From(t *testing.T) {
	testCases := []struct {
		desc string
		prop string
		data []string

		makeExpNode  func(n *movetree.Node)
		expErrSubstr string
	}{
		{
			desc: "black move: pass",
			prop: "B",
			data: []string{},
			makeExpNode: func(n *movetree.Node) {
				n.Properties.Move = move.NewPass(color.Black)
			},
		},
		{
			desc: "black move: pass empty str",
			prop: "B",
			data: []string{""},
			makeExpNode: func(n *movetree.Node) {
				n.Properties.Move = move.NewPass(color.Black)
			},
		},
		{
			desc: "black move",
			prop: "B",
			data: []string{"ab"},
			makeExpNode: func(n *movetree.Node) {
				n.Properties.Move = move.NewMove(color.Black, point.New(0, 1))
			},
		},
		{
			desc: "white move",
			prop: "W",
			data: []string{"ab"},
			makeExpNode: func(n *movetree.Node) {
				n.Properties.Move = move.NewMove(color.White, point.New(0, 1))
			},
		},
		{
			desc: "black placements",
			prop: "AB",
			data: []string{"aa", "bb"},
			makeExpNode: func(n *movetree.Node) {
				n.Properties.Placements = []*move.Move{
					move.NewMove(color.Black, point.New(0, 0)),
					move.NewMove(color.Black, point.New(1, 1)),
				}
			},
		},
		{
			desc: "white placements",
			prop: "AW",
			data: []string{"aa", "bb"},
			makeExpNode: func(n *movetree.Node) {
				n.Properties.Placements = []*move.Move{
					move.NewMove(color.White, point.New(0, 0)),
					move.NewMove(color.White, point.New(1, 1)),
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			n := movetree.NewNode()
			expNode := movetree.NewNode()
			tc.makeExpNode(expNode)
			err := ProcessPropertyData(n, tc.prop, tc.data)
			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Fatal(cerr)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(n, expNode) {
				t.Errorf("got node %v, but expected node %v", n, expNode)
			}
		})
	}
}

func TestConverters_ConvertNode(t *testing.T) {
	testCases := []struct {
		desc     string
		makeNode func(*movetree.Node)

		expOut       string
		expErrSubstr string
	}{
		{
			desc: "black move: pass",
			makeNode: func(n *movetree.Node) {
				n.Properties.Move = move.NewPass(color.Black)
			},
			expOut: "B[]",
		},
		{
			desc: "black move: non-pass",
			makeNode: func(n *movetree.Node) {
				n.Properties.Move = move.NewMove(color.Black, point.New(0, 1))
			},
			expOut: "B[ab]",
		},
		{
			desc: "white move: non-pass",
			makeNode: func(n *movetree.Node) {
				n.Properties.Move = move.NewMove(color.White, point.New(0, 1))
			},
			expOut: "W[ab]",
		},
		{
			desc: "black placements",
			makeNode: func(n *movetree.Node) {
				n.Properties.Placements = []*move.Move{
					move.NewMove(color.Black, point.New(0, 1)),
					move.NewMove(color.Black, point.New(0, 2)),
				}
			},
			expOut: "AB[ab][ac]",
		},
		{
			desc: "white placements",
			makeNode: func(n *movetree.Node) {
				n.Properties.Placements = []*move.Move{
					move.NewMove(color.White, point.New(0, 1)),
					move.NewMove(color.White, point.New(0, 2)),
				}
			},
			expOut: "AW[ab][ac]",
		},
		{
			desc: "black move, extra properties",
			makeNode: func(n *movetree.Node) {
				n.Properties.Move = move.NewMove(color.Black, point.New(0, 1))
				n.SGFProperties["ZZ"] = []string{"zork"}
			},
			expOut: "B[ab]ZZ[zork]",
		},
		{
			desc: "extra properties, sorting",
			makeNode: func(n *movetree.Node) {
				n.Properties.Move = move.NewMove(color.Black, point.New(0, 1))
				n.SGFProperties["ZZ"] = []string{"zork"}
				n.SGFProperties["AA"] = []string{"ark"}
				n.SGFProperties["BB"] = []string{"bark"}
			},
			expOut: "B[ab]AA[ark]BB[bark]ZZ[zork]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			node := movetree.NewNode()
			tc.makeNode(node)
			out, err := ConvertNode(node)
			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Fatal(cerr)
			}
			if err != nil {
				return
			}

			if out != tc.expOut {
				t.Errorf("ConvertNode(%v)=%v, but expected %v", node, out, tc.expOut)
			}
		})
	}
}
