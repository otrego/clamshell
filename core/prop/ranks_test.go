package prop

import (
	"testing"

	"github.com/otrego/clamshell/core/movetree"
)

func TestConvertFromSGF_Rank(t *testing.T) {
	testCases := []fromSGFTestCase{
		{
			desc: "Black Rank",
			prop: "BR",
			data: []string{"20k"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					BlackRank: "20k",
				}
			},
		},
		{
			desc: "White Rank",
			prop: "WR",
			data: []string{"20k"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					WhiteRank: "20k",
				}
			},
		},
		{
			desc: "kyu rank",
			prop: "BR",
			data: []string{"20kyu"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					BlackRank: "20kyu",
				}
			},
		},
		{
			desc: "d rank",
			prop: "BR",
			data: []string{"2d"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					BlackRank: "2d",
				}
			},
		},
		{
			desc: "dan rank",
			prop: "BR",
			data: []string{"2dan"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					BlackRank: "2dan",
				}
			},
		},
		{
			desc: "p rank",
			prop: "BR",
			data: []string{"2p"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					BlackRank: "2p",
				}
			},
		},
		{
			desc: "pro rank",
			prop: "BR",
			data: []string{"2pro"},
			makeExpNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					BlackRank: "2pro",
				}
			},
		},
		{
			desc:         "Invalid rank - extra character",
			prop:         "BR",
			data:         []string{"2prov"},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "Invalid Rank",
		},
		{
			desc:         "Invalid rank - number at end",
			prop:         "BR",
			data:         []string{"2pro2"},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "Invalid Rank",
		},
		{
			desc:         "Invalid rank - no rank",
			prop:         "BR",
			data:         []string{"7"},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "Invalid Rank",
		},
		{
			desc:         "Invalid number - too high",
			prop:         "BR",
			data:         []string{"10p"},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "Invalid number",
		},
		{
			desc:         "Invalid number - too high",
			prop:         "BR",
			data:         []string{"31k"},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "Invalid number",
		},
		{
			desc:         "Invalid number - too low",
			prop:         "BR",
			data:         []string{"0p"},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "Invalid number",
		},
		{
			desc:         "Character in number",
			prop:         "BR",
			data:         []string{"1s3kyu"},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "parsing",
		},
		{
			desc:         "No number",
			prop:         "BR",
			data:         []string{"p"},
			makeExpNode:  func(n *movetree.Node) {},
			expErrSubstr: "parsing",
		},
	}

	testConvertFromSGFCases(t, testCases)
}

func TestConvertNode_Rank(t *testing.T) {
	testCases := []convertNodeTestCase{
		{
			desc: "Both Rank",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					BlackRank: "20k",
					WhiteRank: "7d",
				}
			},
			expOut: "BR[20k]WR[7d]",
		},
		{
			desc: "Black Rank",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					BlackRank: "20k",
				}
			},
			expOut: "BR[20k]",
		},
		{
			desc: "White Rank",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					WhiteRank: "7d",
				}
			},
			expOut: "WR[7d]",
		},
		{
			desc: "Invalid number - too high",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					WhiteRank: "20d",
				}
			},
			expErrSubstr: "Invalid number",
		},
		{
			desc: "Invalid number - too high",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					WhiteRank: "40k",
				}
			},
			expErrSubstr: "Invalid number",
		},
		{
			desc: "Invalid number - too low",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					WhiteRank: "0k",
				}
			},
			expErrSubstr: "Invalid number",
		},
		{
			desc: "Number conversion fail",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					WhiteRank: "20sd",
				}
			},
			expErrSubstr: "parsing",
		},
		{
			desc: "Invalid Rank",
			makeNode: func(n *movetree.Node) {
				n.GameInfo = &movetree.GameInfo{
					WhiteRank: "20dyu",
				}
			},
			expErrSubstr: "Invalid Rank",
		},
	}

	testConvertNodeCases(t, testCases)
}
