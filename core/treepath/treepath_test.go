package treepath

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/otrego/clamshell/core/errcheck"
)

func TestParseInitialPath(t *testing.T) {
	testCases := []struct {
		desc         string
		path         string
		exp          []int
		expErrSubstr string
	}{
		{
			desc: "root move",
			path: "0",
			exp:  []int{},
		},
		{
			desc: "first move",
			path: "1",
			exp:  []int{0},
		},
		{
			desc: "move 10",
			path: "10",
			exp:  []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			desc: "move 2.3",
			path: "2.3",
			exp:  []int{0, 0, 3},
		},
		{
			desc: "move 3",
			path: "3",
			exp:  []int{0, 0, 0},
		},
		{
			desc: "move 2.0 = move 3",
			path: "2.0",
			exp:  []int{0, 0, 0},
		},
		{
			desc: "move 0.0.0.0 = move 3",
			path: "0.0.0.0",
			exp:  []int{0, 0, 0},
		},
		{
			desc: "move 0.0:3 = move 3",
			path: "0.0:3",
			exp:  []int{0, 0, 0},
		},
		// Error cases
		{
			desc:         "move 0:2 is invalid - no repeat on initial pos",
			path:         "0:2",
			expErrSubstr: "unexpected char",
		},
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
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := ParseInitialPath(tc.path)

			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			if !cmp.Equal(got, Treepath(tc.exp)) {
				t.Errorf("ParseInitialPath(%v)=%v, but expected %v", tc.path, got, tc.exp)
			}
		})
	}
}

func TestParseFragment(t *testing.T) {
	testCases := []struct {
		desc         string
		path         string
		exp          []int
		expErrSubstr string
	}{
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
			desc: "varitaion 3",
			path: "3",
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
			got, err := ParseFragment(tc.path)

			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			if !cmp.Equal(got, Treepath(tc.exp)) {
				t.Errorf("ParseInitialPath(%v)=%v, but expected %v", tc.path, got, tc.exp)
			}
		})
	}
}