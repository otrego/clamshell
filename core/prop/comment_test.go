package prop

import (
	"testing"

	"github.com/otrego/clamshell/core/movetree"
)

func TestConvertFromSGF_Comment(t *testing.T) {
	testCases := []fromSGFTestCase{
		{
			desc: "basic comment conversion",
			prop: "C",
			data: []string{"Ima comment"},
			makeExpNode: func(n *movetree.Node) {
				n.Comment = "Ima comment"
			},
		},
	}
	testConvertFromSGFCases(t, testCases)
}

func TestConvertNode_Comment(t *testing.T) {
	testCases := []convertNodeTestCase{
		{
			desc: "basic comment conversion",
			makeNode: func(n *movetree.Node) {
				n.Comment = "Ima comment"
			},
			expOut: "C[Ima comment]",
		},
	}

	testConvertNodeCases(t, testCases)
}
