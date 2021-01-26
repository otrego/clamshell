package prop

import (
	"testing"

	"github.com/otrego/clamshell/core/movetree"
)

func TestConvertFromSGF_Size(t *testing.T) {
	testCases := []fromSGFTestCase{
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
	}

	testConvertFromSGFCases(t, testCases)
}

func TestConvertNode_Size(t *testing.T) {
	testCases := []convertNodeTestCase{
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
	}

	testConvertNodeCases(t, testCases)
}
