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

		expNode      *movetree.Node
		expErrSubstr string
	}{
		{
			desc: "black move: pass",
			prop: "B",
			data: []string{},
			expNode: func(n *movetree.Node) *movetree.Node {
				n.Move = move.NewPass(color.Black)
				return n
			}(movetree.NewNode()),
		},
		{
			desc: "black move: pass v2",
			prop: "B",
			data: []string{""},
			expNode: func(n *movetree.Node) *movetree.Node {
				n.Move = move.NewPass(color.Black)
				return n
			}(movetree.NewNode()),
		},
		{
			desc: "black move",
			prop: "B",
			data: []string{"ab"},
			expNode: func(n *movetree.Node) *movetree.Node {
				n.Move = move.NewMove(color.Black, point.New(0, 1))
				return n
			}(movetree.NewNode()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			n := movetree.NewNode()
			conv := PropConverter(tc.prop)
			if conv == nil {
				t.Fatal("expected converter, but found none")
			}
			if tc.prop != string(conv.Prop) {
				t.Errorf("got converter with prop %q, but expected it to be %q", tc.prop, conv.Prop)
			}
			err := conv.From(n, tc.data)
			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Fatal(cerr)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(n, tc.expNode) {
				t.Errorf("got node %v, but expected node %v", n, tc.expNode)
			}
		})
	}
}

func TestConverters_To(t *testing.T) {
	testCases := []struct {
		desc string
		prop string
		n    *movetree.Node

		expData      []string
		expErrSubstr string
	}{
		{
			desc: "black move: pass",
			prop: "B",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			conv := PropConverter(tc.prop)
			if conv == nil {
				t.Fatal("expected converter, but found none")
			}
			if tc.prop != string(conv.Prop) {
				t.Errorf("got converter with prop %q, but expected it to be %q", tc.prop, conv.Prop)
			}
		})
	}
}
