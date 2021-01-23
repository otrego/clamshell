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
	}

	testConvertNodeCases(t, testCases)
}
