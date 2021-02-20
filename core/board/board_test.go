package board

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			b:    New(9),
			exp:  9,
		},
		{
			desc: "19x19 board",
			b:    New(19),
			exp:  19,
		},
		{
			desc: "13x13 board",
			b:    New(13),
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
			b:    New(9),
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
			b:    New(9),
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
		expCaptures  move.List
		expErrSubstr string
	}{
		{
			desc: "successful added stone -- bottom right",
			b:    New(9),
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
			desc: "successful added stone -- top left",
			b:    New(9),
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
			desc: "successful added stone -- top left capture",
			b: &Board{
				board: [][]color.Color{
					{"", "W", "B", "", "", "", "", "", ""},
					{"", "B", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""}},
			},
			m: move.New(color.Black, point.New(0, 0)),
			exp: "[B * B . . . . . .]\n" +
				"[. B . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]\n" +
				"[. . . . . . . . .]",
			expCaptures: move.List{move.New(color.White, point.New(1, 0))},
		},
		{
			desc: "4 capture groups",
			b: &Board{
				board: [][]color.Color{
					{"", "", "", "", "B", "", "", "", ""},
					{"", "", "", "B", "W", "B", "", "", ""},
					{"", "", "", "B", "W", "B", "", "", ""},
					{"", "B", "B", "B", "W", "B", "B", "B", ""},
					{"B", "W", "W", "W", "", "W", "W", "W", "B"},
					{"", "B", "B", "B", "W", "B", "B", "B", ""},
					{"", "", "", "B", "W", "B", "", "", ""},
					{"", "", "", "B", "W", "B", "", "", ""},
					{"", "", "", "", "B", "", "", "", ""}},
			},
			m: move.New(color.Black, point.New(4, 4)),
			expCaptures: move.List{
				move.New(color.White, point.New(1, 4)),
				move.New(color.White, point.New(2, 4)),
				move.New(color.White, point.New(3, 4)),

				move.New(color.White, point.New(4, 1)),
				move.New(color.White, point.New(4, 2)),
				move.New(color.White, point.New(4, 3)),

				move.New(color.White, point.New(4, 5)),
				move.New(color.White, point.New(4, 6)),
				move.New(color.White, point.New(4, 7)),

				move.New(color.White, point.New(5, 4)),
				move.New(color.White, point.New(6, 4)),
				move.New(color.White, point.New(7, 4)),
			},
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
			b: &Board{
				board: [][]color.Color{{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "B", "", "", "", ""},
					{"", "", "", "B", "", "B", "", "", ""},
					{"", "", "", "B", "B", "W", "", "", ""},
					{"", "", "", "", "W", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""}},
			},
			m:            move.New(color.White, point.New(33, 4)),
			expErrSubstr: "out of bound",
		},
		{
			desc: "test occupied",
			b: &Board{
				board: [][]color.Color{{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "B", "", "", "", ""},
					{"", "", "", "B", "", "B", "", "", ""},
					{"", "", "", "B", "B", "W", "", "", ""},
					{"", "", "", "", "W", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""}},
			},
			m:            move.New(color.White, point.New(4, 3)),
			expErrSubstr: "occupied",
		},
		{
			desc: "test suicidal",
			b: &Board{
				board: [][]color.Color{{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "B", "", "", "", ""},
					{"", "", "", "B", "", "B", "", "", ""},
					{"", "", "", "B", "B", "W", "", "", ""},
					{"", "", "", "", "W", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""}},
			},
			m:            move.New(color.White, point.New(4, 4)),
			expErrSubstr: "suicidal",
		},
		{
			desc: "test ko",
			b: &Board{
				board: [][]color.Color{{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "B", "", "", "", ""},
					{"", "", "", "B", "", "B", "", "", ""},
					{"", "", "", "W", "B", "W", "", "", ""},
					{"", "", "", "", "W", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""},
					{"", "", "", "", "", "", "", "", ""}},
				ko: point.New(4, 4),
			},
			m:            move.New(color.White, point.New(4, 4)),
			expErrSubstr: "illegal",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			capt, err := tc.b.PlaceStone(tc.m)
			got := tc.b.String()

			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(capt, tc.expCaptures) {
				t.Errorf("Got captures:%v, but expected captures:%v; diff=%v", capt, tc.expCaptures,
					cmp.Diff(capt.String(), tc.expCaptures.String()))
			}

			if got != tc.exp {
				t.Errorf("got board:\n%v, but expected board:\n%v", got, tc.exp)
			}
		})
	}
}

func TestClone(t *testing.T) {
	b := &Board{
		board: [][]color.Color{{"", "", "", "", "", "", "", "", ""},
			{"", "", "", "", "", "", "", "", ""},
			{"", "", "", "", "", "", "", "", ""},
			{"", "", "", "", "B", "", "", "", ""},
			{"", "", "", "B", "", "B", "", "", ""},
			{"", "", "", "W", "B", "W", "", "", ""},
			{"", "", "", "", "W", "", "", "", ""},
			{"", "", "", "", "", "", "", "", ""},
			{"", "", "", "", "", "", "", "", ""}},
		ko: point.New(4, 5),
	}

	newb := b.Clone()

	if !reflect.DeepEqual(b, newb) {
		t.Fatalf("b.Clone()=%v, but expected %v. diff=%v", b.String(), newb.String(), cmp.Diff(b.String(), newb.String()))
	}
}

func TestStoneState(t *testing.T) {
	b := New(5)
	err := b.SetPlacements(move.List{
		move.New(color.Black, point.New(2, 3)),
		move.New(color.White, point.New(2, 4)),
		move.New(color.White, point.New(1, 1)),
	})
	if err != nil {
		t.Fatal(err)
	}

	fbstate := b.StoneState()
	expState := move.List{
		move.New(color.White, point.New(1, 1)),
		move.New(color.Black, point.New(2, 3)),
		move.New(color.White, point.New(2, 4)),
	}
	if !reflect.DeepEqual(fbstate, expState) {
		t.Errorf("got stone state %v, but expected %v", fbstate, expState)
	}
}

func TestFullBoardState(t *testing.T) {
	b := New(5)
	err := b.SetPlacements(move.List{
		move.New(color.Black, point.New(2, 3)),
		move.New(color.White, point.New(2, 4)),
		move.New(color.White, point.New(1, 1)),
	})
	if err != nil {
		t.Fatal(err)
	}

	fbstate := b.FullBoardState()
	expState := [][]color.Color{
		[]color.Color{color.Empty, color.Empty, color.Empty, color.Empty, color.Empty},
		[]color.Color{color.Empty, color.White, color.Empty, color.Empty, color.Empty},
		[]color.Color{color.Empty, color.Empty, color.Empty, color.Empty, color.Empty},
		[]color.Color{color.Empty, color.Empty, color.Black, color.Empty, color.Empty},
		[]color.Color{color.Empty, color.Empty, color.White, color.Empty, color.Empty},
	}
	if !cmp.Equal(fbstate, expState) {
		t.Errorf("got full stone state %v, but expected %v. diff=%v", fbstate, expState, cmp.Diff(fbstate, expState))
	}
}
