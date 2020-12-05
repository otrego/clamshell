package movetree_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/otrego/clamshell/core/errcheck"
	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/sgf"
)

func TestParsePath(t *testing.T) {
	testCases := []struct {
		desc         string
		path         string
		exp          []int
		expErrSubstr string
	}{
		{
			desc: "empty treepath",
			path: "",
			exp:  []int{},
		},
		{
			desc: "empty treepath, using leading .",
			path: ".",
			exp:  []int{},
		},
		{
			desc: "0th varation",
			path: "0",
			exp:  []int{0},
		},
		{
			desc: "first move",
			path: "1",
			exp:  []int{1},
		},
		{
			desc: "variation 10",
			path: "10",
			exp:  []int{10},
		},
		{
			desc: "move 2.3",
			path: "2.3",
			exp:  []int{2, 3},
		},
		{
			desc: "variation 3",
			path: "3",
			exp:  []int{3},
		},
		{
			desc: "variation 3, with leading .",
			path: ".3",
			exp:  []int{3},
		},
		{
			desc: "move 2.0",
			path: "2.0",
			exp:  []int{2, 0},
		},
		{
			desc: "move 0.0.0.0 = 4 moves",
			path: "0.0.0.0",
			exp:  []int{0, 0, 0, 0},
		},
		{
			desc: "move 0.0.0.0 = 4 moves, leading .",
			path: ".0.0.0.0",
			exp:  []int{0, 0, 0, 0},
		},
		{
			desc: "move 0.0:3 = 4 moves",
			path: "0.0:3",
			exp:  []int{0, 0, 0, 0},
		},
		{
			desc: "move 0:4 = 4 moves",
			path: "0:4",
			exp:  []int{0, 0, 0, 0},
		},
		{
			desc: "move 0:4 = 4 moves @ var 1",
			path: "1:4",
			exp:  []int{1, 1, 1, 1},
		},
		{
			desc: "complicated pattern",
			path: "1.2:1.0.3:3",
			exp:  []int{1, 2, 0, 3, 3, 3},
		},
		// Error cases
		{
			desc:         "move a:2 is invalid - not a number",
			path:         "a:2",
			expErrSubstr: "unexpected char",
		},
		{
			desc:         "move -1:2 is invalid - negatives are not allowed",
			path:         "-1:2",
			expErrSubstr: "unexpected char",
		},
		{
			desc:         "move 1::2 is invalid - no repeat separators",
			path:         "1::2",
			expErrSubstr: "separator",
		},
		{
			desc:         "move 1..2 is invalid - no repeat separators",
			path:         "1..2",
			expErrSubstr: "separator",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := movetree.ParsePath(tc.path)

			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			if !cmp.Equal(got, movetree.Treepath(tc.exp)) {
				t.Errorf("ParsePath(%v)=%v, but expected %v", tc.path, got, tc.exp)
			}
		})
	}
}

func TestApplyPath(t *testing.T) {
	testCases := []struct {
		desc     string
		path     string
		game     string
		expProps map[string][]string
	}{
		{
			desc: "first move",
			path: ".0",
			game: "(;GM[1];PM[1]B[pd]C[foo])",
			expProps: map[string][]string{
				"C":  []string{"foo"},
				"PM": []string{"1"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := sgf.FromString(tc.game).Parse()
			if err != nil {
				t.Error(err)
				return
			}
			path, err := movetree.ParsePath(tc.path)
			if err != nil {
				t.Error(err)
				return
			}
			n := path.Apply(g.Root)

			if !cmp.Equal(n.Properties, tc.expProps) {
				t.Errorf("path.Apply(root)=%v, expected %v. Diff=%v", n.Properties, tc.expProps, cmp.Diff(n.Properties, tc.expProps))
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		desc string
		tp   movetree.Treepath
		exp  string
	}{
		{
			desc: "empty treepath",
			tp:   []int{},
			exp:  "[]",
		},
		{
			desc: "short treepath",
			tp:   []int{1},
			exp:  "[1]",
		},
		{
			desc: "long treepath",
			tp:   []int{1, 2, 0, 2, 2, 2},
			exp:  "[1 2 0 2 2 2]",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.tp.String()
			if got != tc.exp {
				t.Errorf("got %s, expected %s", got, tc.exp)
			}
		})
	}
}

func TestCompactString(t *testing.T) {
	testCases := []struct {
		desc string
		tp   movetree.Treepath
		exp  string
	}{
		{
			desc: "empty treepath",
			tp:   []int{},
			exp:  ".",
		},
		{
			desc: "short treepath",
			tp:   []int{1},
			exp:  ".1",
		},
		{
			desc: "long repeat",
			tp:   []int{1, 1, 1, 1},
			exp:  ".1:4",
		},
		{
			desc: "long no repeat",
			tp:   []int{1, 2, 3, 4, 5, 6},
			exp:  ".1.2.3.4.5.6",
		},
		{
			desc: "long mixed",
			tp:   []int{1, 2, 2, 0, 0, 5, 0, 2, 2, 2, 2, 1},
			exp:  ".1.2:2.0:2.5.0.2:4.1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.tp.CompactString()
			if got != tc.exp {
				t.Errorf("got %s, expected %s", got, tc.exp)
			}
		})
	}
}
