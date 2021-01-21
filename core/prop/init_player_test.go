package prop

import (
	"testing"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/movetree"
)

func TestConvertFromSGF_InitPlayer(t *testing.T) {
	testCases := []fromSGFTestCase{
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

	testConvertFromSGFCases(t, testCases)
}

func TestConvertNode_InitPlayer(t *testing.T) {
	testCases := []convertNodeTestCase{
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

	testConvertNodeCases(t, testCases)
}
