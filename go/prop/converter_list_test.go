package prop

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/move"
	"github.com/otrego/clamshell/go/movetree"
	"github.com/otrego/clamshell/go/point"
)

type propGetter func(*movetree.Node) interface{}

type setprops func(*movetree.Node) *movetree.Node

type fromSGFTestCase struct {
	desc        string
	prop        string
	data        []string
	makeExpNode func(*movetree.Node)
	expErr      error
}

func testConvertFromSGFCases(t *testing.T, testCases []fromSGFTestCase) {
	for _, tci := range testCases {
		tc := tci
		t.Run(tc.desc, func(t *testing.T) {
			n := movetree.NewNode()
			expNode := movetree.NewNode()
			tc.makeExpNode(expNode)

			err := ProcessPropertyData(n, tc.prop, tc.data)
			if !errors.Is(err, tc.expErr) {
				t.Fatalf("got error %v, but expected err %v", err, tc.expErr)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(n, expNode) {
				t.Errorf("got node %#v, but expected node %#v", n, expNode)
			}
		})
	}
}

type convertNodeTestCase struct {
	desc     string
	makeNode func(*movetree.Node)
	expOut   string
	expErr   error
}

func testConvertNodeCases(t *testing.T, testCases []convertNodeTestCase) {
	for _, tci := range testCases {
		tc := tci
		t.Run(tc.desc, func(t *testing.T) {
			node := movetree.NewNode()
			tc.makeNode(node)
			out, err := ConvertNode(node)
			if !errors.Is(err, tc.expErr) {
				t.Fatalf("got error %v, but expected %v", err, tc.expErr)
			}
			if err != nil {
				return
			}

			if out != tc.expOut {
				diff := cmp.Diff(out, tc.expOut)
				t.Errorf("ConvertNode(%#v)=%v, but expected %v. Diff=%s", node, out, tc.expOut, diff)
			}
		})
	}
}

func TestConvertFromSGF_Collection(t *testing.T) {
	testCases := []fromSGFTestCase{
		{
			desc: "unknown property",
			prop: "ZZ",
			data: []string{"Zork"},
			makeExpNode: func(n *movetree.Node) {
				n.SGFProperties["ZZ"] = []string{"Zork"}
			},
		},
	}

	testConvertFromSGFCases(t, testCases)
}

func TestConvertNode_Collection(t *testing.T) {
	testCases := []convertNodeTestCase{
		{
			desc: "black move, extra properties",
			makeNode: func(n *movetree.Node) {
				n.Move = move.New(color.Black, point.New(0, 1))
				n.SGFProperties["ZZ"] = []string{"zork"}
			},
			expOut: "B[ab]ZZ[zork]",
		},
		{
			desc: "extra properties, sorting",
			makeNode: func(n *movetree.Node) {
				n.Move = move.New(color.Black, point.New(0, 1))
				n.SGFProperties["ZZ"] = []string{"zork"}
				n.SGFProperties["AA"] = []string{"ark"}
				n.SGFProperties["BB"] = []string{"bark"}
			},
			expOut: "B[ab]AA[ark]BB[bark]ZZ[zork]",
		},
	}

	testConvertNodeCases(t, testCases)
}
