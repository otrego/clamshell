package board

import (
	"fmt"
	"testing"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/errcheck"
	"github.com/otrego/clamshell/core/move"
	"github.com/otrego/clamshell/core/point"
)

func TestNewBoard(t *testing.T) {
	testCases := []struct {
		desc string
		b    *Board
		exp  int
	}{
		{
			desc: "9x9 board",
			b:    NewBoard(9),
			exp:  9,
		},
		{
			desc: "19x19 board",
			b:    NewBoard(19),
			exp:  19,
		},
		{
			desc: "13x13 board",
			b:    NewBoard(13),
			exp:  13,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got1 := len(tc.b.board)
			got2 := len(tc.b.board[0])
			if got1 != tc.exp || got2 != tc.exp {
				t.Errorf("got %dx%d, expected %v", got1, got2, tc.exp)
			}
		})
	}
}

func TestString(t *testing.T) {

	testCases := []struct {
		desc string
		b    *Board
		exp  string
	}{
		{
			desc: "empty 9x9 board",
			b:    NewBoard(9),
			exp: "[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]",
		},
		{
			desc: "some White and Black added 9x9 board",
			b: &Board{[][]color.Color{{"", "", "", "", "", "", "B", "W", ""},
				{"B", "", "", "", "", "B", "W", "W", ""},
				{"B", "", "", "W", "", "", "B", "W", ""},
				{"W", "", "", "", "", "B", "B", "W", ""},
				{"", "", "", "B", "W", "B", "W", "B", "W"},
				{"", "B", "", "", "", "", "W", "B", "B"},
				{"", "", "W", "", "B", "", "W", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
				nil,
			},
			exp: "[. . . . . . B W .]\n" +
				"[B . . . . B W W .]\n" +
				"[B . . W . . B W .]\n" +
				"[W . . . . B B W .]\n" +
				"[. . . B W B W B W]\n" +
				"[. B . . . . W B B]\n" +
				"[. . W . B . W . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			got := tc.b.String()
			if got != tc.exp {
				t.Errorf("got %s, exptected %s", got, tc.exp)
			}
		})
	}
}

func TestCapturedStones(t *testing.T) {
	testCases := []struct {
		desc string
		b    *Board
		pt   *point.Point
		exp  []*point.Point
	}{
		{
			desc: "empty board",
			b:    NewBoard(9),
			pt:   point.New(5, 5),
			exp:  nil,
		},
		{
			desc: "some White and Black added 9x9 board, no captures",
			b: &Board{[][]color.Color{{"", "", "", "", "", "", "B", "W", ""},
				{"B", "", "", "", "", "B", "W", "W", ""},
				{"B", "", "", "W", "", "", "B", "W", ""},
				{"W", "", "", "", "", "B", "B", "W", ""},
				{"", "", "", "B", "W", "B", "W", "B", "W"},
				{"", "B", "", "", "", "", "W", "B", "B"},
				{"", "", "W", "", "B", "", "W", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
				nil,
			},
			pt:  point.New(5, 5),
			exp: nil,
		},
		{
			desc: "deep liberty",
			b: &Board{[][]color.Color{{"", "", "", "", "", "", "", "", ""},
				{"", "B", "B", "B", "B", "B", "B", "", ""},
				{"", "B", "W", "W", "W", "W", "W", "", ""},
				{"", "B", "W", "B", "B", "B", "B", "", ""},
				{"", "B", "W", "B", "W", "W", "B", "", ""},
				{"", "B", "W", "B", "B", "W", "B", "", ""},
				{"", "B", "W", "W", "W", "W", "B", "", ""},
				{"", "B", "B", "B", "B", "B", "B", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
				nil,
			},
			pt:  point.New(4, 4),
			exp: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.b.capturedStones(tc.pt)
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tc.exp) {
				t.Errorf("got %v, but expected %v", got, tc.exp)
			}
		})
	}
}

func TestRemoveCapturedStones(t *testing.T) {
	testCases := []struct {
		desc string
		b    *Board
		m    *move.Move
		exp  string
	}{
		{
			desc: "4 captures",
			b: &Board{[][]color.Color{{"", "", "", "", "B", "", "", "", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "B", "B", "B", "W", "B", "B", "B", ""},
				{"B", "W", "W", "W", "B", "W", "W", "W", "B"},
				{"", "B", "B", "B", "W", "B", "B", "B", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "", "", "", "B", "", "", "", ""}},
				nil,
			},
			m: move.New(color.Black, point.New(4, 4)),
			exp: "[. . . . B . . . .]\n" +
				"[. . . B . B . . .]\n" +
				"[. . . B . B . . .]\n" +
				"[. B B B . B B B .]\n" +
				"[B . . . B . . . B]\n" +
				"[. B B B . B B B .]\n" +
				"[. . . B . B . . .]\n" +
				"[. . . B . B . . .]\n" +
				"[. . . . B . . . .]",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			capturedStones := tc.b.findCapturedGroups(tc.m)
			tc.b.removeCapturedStones(capturedStones)
			got := tc.b.String()
			if fmt.Sprintf("%s", got) != fmt.Sprintf("%s", tc.exp) {
				t.Errorf("got %v, but expected %v", got, tc.exp)
			}
		})
	}
}

func TestPlaceStone(t *testing.T) {
	testCases := []struct {
		desc         string
		b            *Board
		m            *move.Move
		exp          string
		expErrSubstr string
	}{
		{
			desc: "successful added stone",
			b:    NewBoard(9),
			m:    move.New(color.Black, point.New(8, 8)),
			exp: "[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . B]",
		},
		{
			desc: "successful added stone",
			b:    NewBoard(9),
			m:    move.New(color.Black, point.New(0, 0)),
			exp: "[B . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]",
		},
		{
			desc: "4 captures",
			b: &Board{[][]color.Color{{"", "", "", "", "B", "", "", "", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "B", "B", "B", "W", "B", "B", "B", ""},
				{"B", "W", "W", "W", "", "W", "W", "W", "B"},
				{"", "B", "B", "B", "W", "B", "B", "B", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "", "", "", "B", "", "", "", ""}},
				nil,
			},
			m: move.New(color.Black, point.New(4, 4)),
			exp: "[. . . . B . . . .]\n" +
				"[. . . B . B . . .]\n" +
				"[. . . B . B . . .]\n" +
				"[. B B B . B B B .]\n" +
				"[B . . . B . . . B]\n" +
				"[. B B B . B B B .]\n" +
				"[. . . B . B . . .]\n" +
				"[. . . B . B . . .]\n" +
				"[. . . . B . . . .]",
		},
		{
			desc: "test out of bounds",
			b: &Board{[][]color.Color{{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "B", "", "", "", ""},
				{"", "", "", "B", "", "B", "", "", ""},
				{"", "", "", "B", "B", "W", "", "", ""},
				{"", "", "", "", "W", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
				nil,
			},
			m:            move.New(color.White, point.New(33, 4)),
			expErrSubstr: "out of bound",
		},
		{
			desc: "test occupied",
			b: &Board{[][]color.Color{{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "B", "", "", "", ""},
				{"", "", "", "B", "", "B", "", "", ""},
				{"", "", "", "B", "B", "W", "", "", ""},
				{"", "", "", "", "W", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
				nil,
			},
			m:            move.New(color.White, point.New(4, 3)),
			expErrSubstr: "occupied",
		},
		{
			desc: "test suicidal",
			b: &Board{[][]color.Color{{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "B", "", "", "", ""},
				{"", "", "", "B", "", "B", "", "", ""},
				{"", "", "", "B", "B", "W", "", "", ""},
				{"", "", "", "", "W", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
				nil,
			},
			m:            move.New(color.White, point.New(4, 4)),
			expErrSubstr: "suicidal",
		},
		{
			desc: "test ko",
			b: &Board{[][]color.Color{{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "B", "", "", "", ""},
				{"", "", "", "B", "", "B", "", "", ""},
				{"", "", "", "W", "B", "W", "", "", ""},
				{"", "", "", "", "W", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
				point.New(4, 5),
			},
			m:            move.New(color.White, point.New(4, 4)),
			expErrSubstr: "illegal",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := tc.b.PlaceStone(tc.m)
			got := tc.b.String()

			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			if fmt.Sprintf("%s", got) != fmt.Sprintf("%s", tc.exp) {
				t.Errorf("got %v, but expected %v", got, tc.exp)
			}
		})
	}
}
