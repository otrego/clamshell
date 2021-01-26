package prop

import (
	"testing"

	"github.com/otrego/clamshell/core/movetree"
)

func TestConvertFromSGF_Komi(t *testing.T) {
	testCases := []fromSGFTestCase{
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
	}

	testConvertFromSGFCases(t, testCases)
}

func TestConvertNode_Komi(t *testing.T) {
	testCases := []convertNodeTestCase{
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
	}

	testConvertNodeCases(t, testCases)
}
