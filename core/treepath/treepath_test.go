package treepath

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/otrego/clamshell/core/errcheck"
	"github.com/otrego/clamshell/core/sgf"
)

func TestParse(t *testing.T) {
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
			got, err := Parse(tc.path)

			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			if !cmp.Equal(got, Treepath(tc.exp)) {
				t.Errorf("Parse(%v)=%v, but expected %v", tc.path, got, tc.exp)
			}
		})
	}
}

func TestApplyPath(t *testing.T) {
	testCases := []struct {
		desc     string
		initPath string
		game     string
		expProps map[string][]string
	}{
		{
			desc:     "first move",
			initPath: "0",
			game:     "(;GM[1];B[pd]C[foo])",
			expProps: map[string][]string{
				"C": []string{"zed"},
				"B": []string{"pd"},
			},
		},
	}
	t.Skip("Parsing is currently incorrect:Re-enable test once https://github.com/otrego/clamshell/issues/83 is fixed")
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := sgf.FromString(tc.game).Parse()
			if err != nil {
				t.Error(err)
				return
			}
			path, err := Parse(tc.initPath)
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
