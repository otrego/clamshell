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
			desc: "black placements",
			prop: "AB",
			data: []string{"aa", "bb"},
			makeExpNode: func(n *movetree.Node) {
				n.Placements = []*move.Move{
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
				n.Placements = []*move.Move{
					move.NewMove(color.White, point.New(0, 0)),
					move.NewMove(color.White, point.New(1, 1)),
				}
			},
		},
		{
			desc: "size",
			prop: "SZ",
			data: []string{"13"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Size: 13,
				}
			},
		},
		{
			desc: "Komi",
			prop: "KM",
			data: []string{"3.5"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Komi: new(float64),
				}
				*n.GameInfo.Komi = 3.5
			},
		},
		{
			desc:         "Bad Komi",
			prop:         "KM",
			data:         []string{"3.25"},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "for prop KM",
		},
		{
			desc: "Initial Turn 1",
			prop: "PL",
			data: []string{"B"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Player: color.Black,
				}
			},
		},
		{
			desc: "Initial Turn 2",
			prop: "PL",
			data: []string{"b"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Player: color.Black,
				}
			},
		},
		{
			desc: "Initial Turn 3",
			prop: "PL",
			data: []string{"W"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Player: color.White,
				}
			},
		},
		{
			desc: "Initial Turn 4",
			prop: "PL",
			data: []string{"w"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Player: color.White,
				}
			},
		},
		{
			desc:         "Bad Initial Turn 1",
			prop:         "PL",
			data:         []string{},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "requires exactly 1",
		},
		{
			desc:         "Bad Initial Turn 2",
			prop:         "PL",
			data:         []string{"f"},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "PL has invalid value",
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
			desc: "black placements",
			makeNode: func(n *movetree.Node) {
				n.Placements = []*move.Move{
					move.NewMove(color.Black, point.New(0, 1)),
					move.NewMove(color.Black, point.New(0, 2)),
				}
			},
			expOut: "AB[ab][ac]",
		},
		{
			desc: "white placements",
			makeNode: func(n *movetree.Node) {
				n.Placements = []*move.Move{
					move.NewMove(color.White, point.New(0, 1)),
					move.NewMove(color.White, point.New(0, 2)),
				}
			},
			expOut: "AW[ab][ac]",
		},
		{
			desc: "black move, extra properties",
			makeNode: func(n *movetree.Node) {
				n.Move = move.NewMove(color.Black, point.New(0, 1))
				n.SGFProperties["ZZ"] = []string{"zork"}
			},
			expOut: "B[ab]ZZ[zork]",
		},
		{
			desc: "extra properties, sorting",
			makeNode: func(n *movetree.Node) {
				n.Move = move.NewMove(color.Black, point.New(0, 1))
				n.SGFProperties["ZZ"] = []string{"zork"}
				n.SGFProperties["AA"] = []string{"ark"}
				n.SGFProperties["BB"] = []string{"bark"}
			},
			expOut: "B[ab]AA[ark]BB[bark]ZZ[zork]",
		},
		{
			desc: "size",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Size: 13,
				}
			},
			expOut: "SZ[13]",
		},
		{
			desc: "size, empty",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{}
			},
			expOut: "",
		},
		{
			desc: "size, invalid",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Size: 100,
				}
			},
			expErrSubstr: "invalid board size",
		},
		{
			desc: "komi",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Komi: new(float64),
				}
				*n.GameInfo.Komi = 3.5
			},
			expOut: "KM[3.5]",
		},
		{
			desc: "komi, nil",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{}
			},
			expOut: "",
		},
		{
			desc: "komi, gameInfo nil",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = nil
			},
			expOut: "",
		},
		{
			desc: "komi, gameInfo nil",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = nil
			},
			expOut: "",
		},
		{
			desc: "komi, invalid",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Komi: new(float64),
				}
				*n.GameInfo.Komi = 3.25
			},
			expErrSubstr: "invalid komi",
		},
		{
			desc: "Initial Turn 1",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Player: color.Black,
				}
			},
			expOut: "PL[B]",
		},
		{
			desc: "Initial Turn 2",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Player: color.White,
				}
			},
			expOut: "PL[W]",
		},
		{
			desc: "Initial Turn, gameInfo nil",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = nil
			},
			expOut: "",
		},
		{
			desc: "Initial Turn, gameInfo.Player invalid",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					Player: "g",
				}
			},
			expErrSubstr: "W or B, but was",
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
