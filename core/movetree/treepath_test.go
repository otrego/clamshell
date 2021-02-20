package movetree_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/otrego/clamshell/core/board"
	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/errcheck"
	"github.com/otrego/clamshell/core/move"
	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/point"
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
			desc: "empty treepath, using leading -",
			path: "-",
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
			desc: "move 2-3",
			path: "2-3",
			exp:  []int{2, 3},
		},
		{
			desc: "variation 3",
			path: "3",
			exp:  []int{3},
		},
		{
			desc: "variation 3, with leading -",
			path: "-3",
			exp:  []int{3},
		},
		{
			desc: "move 2-0",
			path: "2-0",
			exp:  []int{2, 0},
		},
		{
			desc: "move 0-0-0-0 = 4 moves",
			path: "0-0-0-0",
			exp:  []int{0, 0, 0, 0},
		},
		{
			desc: "move 0-0-0-0 = 4 moves, leading -",
			path: "-0-0-0-0",
			exp:  []int{0, 0, 0, 0},
		},
		{
			desc: "move 0-0x3 = 4 moves",
			path: "0-0x3",
			exp:  []int{0, 0, 0, 0},
		},
		{
			desc: "move 0x4 = 4 moves",
			path: "0x4",
			exp:  []int{0, 0, 0, 0},
		},
		{
			desc: "move 0x4 = 4 moves @ var 1",
			path: "1x4",
			exp:  []int{1, 1, 1, 1},
		},
		{
			desc: "complicated pattern",
			path: "1-2x1-0-3x3",
			exp:  []int{1, 2, 0, 3, 3, 3},
		},
		// Error cases
		{
			desc:         "move ax2 is invalid - not a number",
			path:         "ax2",
			expErrSubstr: "unexpected char",
		},
		{
			desc:         "move .1x2 is invalid - periods are not allowed",
			path:         ".1x2",
			expErrSubstr: "unexpected char",
		},
		{
			desc:         "move 1xx2 is invalid - no repeat separators",
			path:         "1xx2",
			expErrSubstr: "separator",
		},
		{
			desc:         "move 1--2 is invalid - no repeat separators",
			path:         "1--2",
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

			if !cmp.Equal(got, movetree.Path(tc.exp)) {
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
			path: "-0",
			game: "(;GM[1];PM[1]B[pd]ZZ[foo])",
			expProps: map[string][]string{
				"ZZ": []string{"foo"},
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

			if !cmp.Equal(n.SGFProperties, tc.expProps) {
				t.Errorf("path.Apply(root)=%v, expected %v. Diff=%v", n.SGFProperties, tc.expProps, cmp.Diff(n.SGFProperties, tc.expProps))
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		desc string
		tp   movetree.Path
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
		tp   movetree.Path
		exp  string
	}{
		{
			desc: "empty treepath",
			tp:   []int{},
			exp:  "-",
		},
		{
			desc: "short treepath",
			tp:   []int{1},
			exp:  "-1",
		},
		{
			desc: "long repeat",
			tp:   []int{1, 1, 1, 1},
			exp:  "-1x4",
		},
		{
			desc: "long no repeat",
			tp:   []int{1, 2, 3, 4, 5, 6},
			exp:  "-1-2-3-4-5-6",
		},
		{
			desc: "long mixed",
			tp:   []int{1, 2, 2, 0, 0, 5, 0, 2, 2, 2, 2, 1},
			exp:  "-1-2x2-0x2-5-0-2x4-1",
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
func makeBoard(ml move.List) *board.Board {
	b := board.New(9)
	for _, mv := range ml {
		_, err := b.PlaceStone(mv)
		if err != nil {
			panic(err)
		}
	}
	return b
}

func TestApplyToBoard(t *testing.T) {
	testCases := []struct {
		desc string
		sgf  string
		tp   string
		b    *board.Board

		expBoard     *board.Board
		expCaptures  move.List
		expErrSubstr string
	}{
		{
			desc:     "empty apply case",
			sgf:      "(;GM[1];B[aa];W[ab];B[ac];W[dd];B[bb])",
			b:        makeBoard(move.List{}),
			expBoard: makeBoard(move.List{}),
		},
		{
			desc: "basic apply case",
			sgf:  "(;GM[1];B[aa];W[ab];B[ac];W[dd];B[bb])",
			b:    makeBoard(move.List{}),
			tp:   "0",
			expBoard: makeBoard(move.List{
				move.New(color.Black, point.New(0, 0)),
			}),
		},
		{
			desc: "apply with captures",
			sgf:  "(;GM[1];B[aa];W[ab];B[ac];W[dd];B[bb])",
			b:    makeBoard(move.List{}),
			tp:   "0x5",
			expBoard: makeBoard(move.List{
				move.New(color.Black, point.New(0, 0)),
				move.New(color.Black, point.New(0, 2)),
				move.New(color.White, point.New(0, 1)),
				move.New(color.White, point.New(3, 3)),
				move.New(color.Black, point.New(1, 1)),
			}),
			expCaptures: move.List{
				move.New(color.White, point.New(0, 1)),
			},
		},
		{
			desc: "apply with captures & Placements",
			sgf:  "(;GM[1]AB[ad][ae]AW[bd][be];B[aa];W[ab];B[ac];W[dd];B[bb])",
			b:    makeBoard(move.List{}),
			tp:   "0x5",
			expBoard: makeBoard(move.List{
				move.New(color.Black, point.New(0, 0)),
				move.New(color.Black, point.New(0, 2)),
				move.New(color.Black, point.New(0, 3)),
				move.New(color.Black, point.New(0, 4)),

				move.New(color.White, point.New(0, 1)),
				move.New(color.White, point.New(3, 3)),
				move.New(color.White, point.New(1, 3)),
				move.New(color.White, point.New(1, 4)),

				move.New(color.Black, point.New(1, 1)),
			}),
			expCaptures: move.List{
				move.New(color.White, point.New(0, 1)),
			},
		},
	}

	for _, tci := range testCases {
		tc := tci
		t.Run(tc.desc, func(t *testing.T) {
			mt, err := sgf.FromString(tc.sgf).Parse()
			if err != nil {
				t.Fatal(err)
			}
			tp, err := movetree.ParsePath(tc.tp)
			if err != nil {
				t.Fatal(err)
			}

			b, capt, err := tp.ApplyToBoard(mt.Root, tc.b)
			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
			}
			if err != nil {
				return
			}

			if tc.expBoard == nil {
				t.Fatal("expected board was nil")
			}
			if !reflect.DeepEqual(b, tc.expBoard) {
				t.Errorf("After apply, got board:\n%v, but expected:\n%v. Diff=%v", b, tc.expBoard, cmp.Diff(b.String(), tc.expBoard.String()))
			}

			if !reflect.DeepEqual(capt, tc.expCaptures) {
				t.Errorf("After apply, got captures %v, but expected %v. Diff=%v", capt, tc.expCaptures, cmp.Diff(capt.String(), tc.expCaptures.String()))
			}
		})
	}
}
