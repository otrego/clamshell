package game

import (
	"fmt"
	"testing"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/point"
)

func TestNewBoard(t *testing.T) {
	testCases := []struct {
		desc string
		b    Board
		exp  int
	}{
		{
			desc: "9x9 board",
			b:    *NewBoard(9),
			exp:  9,
		},
		{
			desc: "19x19 board",
			b:    *NewBoard(19),
			exp:  19,
		},
		{
			desc: "13x13 board",
			b:    *NewBoard(13),
			exp:  13,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got1 := len(tc.b)
			got2 := len(tc.b[0])
			if got1 != tc.exp || got2 != tc.exp {
				t.Errorf("got %dx%d, expected %v", got1, got2, tc.exp)
			}
		})
	}
}

func TestString(t *testing.T) {

	testCases := []struct {
		desc string
		b    Board
		exp  string
	}{
		{
			desc: "empty 9x9 board",
			b:    *NewBoard(9),
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
			b: [][]color.Color{{"", "", "", "", "", "", "B", "W", ""},
				{"B", "", "", "", "", "B", "W", "W", ""},
				{"B", "", "", "W", "", "", "B", "W", ""},
				{"W", "", "", "", "", "B", "B", "W", ""},
				{"", "", "", "B", "W", "B", "W", "B", "W"},
				{"", "B", "", "", "", "", "W", "B", "B"},
				{"", "", "W", "", "B", "", "W", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
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

func TestIsCaptured(t *testing.T) {
	testCases := []struct {
		desc string
		b    Board
		pt   *point.Point
		exp  []*point.Point
	}{
		{
			desc: "empty board",
			b:    *NewBoard(9),
			pt:   point.New(5, 5),
			exp:  nil,
		},
		{
			desc: "some White and Black added 9x9 board, no captures",
			b: [][]color.Color{{"", "", "", "", "", "", "B", "W", ""},
				{"B", "", "", "", "", "B", "W", "W", ""},
				{"B", "", "", "W", "", "", "B", "W", ""},
				{"W", "", "", "", "", "B", "B", "W", ""},
				{"", "", "", "B", "W", "B", "W", "B", "W"},
				{"", "B", "", "", "", "", "W", "B", "B"},
				{"", "", "W", "", "B", "", "W", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
			pt:  point.New(5, 5),
			exp: nil,
		},
		{
			desc: "some White and Black added 9x9 board, 2 captures",
			b: [][]color.Color{{"W", "W", "W", "", "", "", "B", "W", ""},
				{"B", "W", "", "W", "", "B", "W", "W", ""},
				{"B", "W", "W", "W", "", "", "B", "W", ""},
				{"W", "W", "", "", "", "B", "B", "W", ""},
				{"", "", "", "B", "W", "B", "W", "B", "W"},
				{"", "B", "", "", "", "", "W", "B", "B"},
				{"", "", "W", "", "B", "", "W", "", ""},
				{"", "", "", "", "", "", "", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
			pt: point.New(0, 1),
			exp: []*point.Point{
				point.New(0, 2),
				point.New(0, 1),
			},
		},
		{
			desc: "deep liberty",
			b: [][]color.Color{{"", "", "", "", "", "", "", "", ""},
				{"", "B", "B", "B", "B", "B", "B", "", ""},
				{"", "B", "W", "W", "W", "W", "W", "", ""},
				{"", "B", "W", "B", "B", "B", "B", "", ""},
				{"", "B", "W", "B", "W", "W", "B", "", ""},
				{"", "B", "W", "B", "B", "W", "B", "", ""},
				{"", "B", "W", "W", "W", "W", "B", "", ""},
				{"", "B", "B", "B", "B", "B", "B", "", ""},
				{"", "", "", "", "", "", "", "", ""}},
			pt:  point.New(4, 4),
			exp: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			_, got := tc.b.IsCaptured(tc.pt)
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tc.exp) {
				t.Errorf("got %v, but expected %v", got, tc.exp)
			}
		})
	}
}

func TestFindEnemyGroup(t *testing.T) {
	testCases := []struct {
		desc string
		b    Board
		pt   *point.Point
		c    color.Color
		exp  string
	}{
		{
			desc: "4 captures",
			b: [][]color.Color{{"", "", "", "", "B", "", "", "", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "B", "B", "B", "W", "B", "B", "B", ""},
				{"B", "W", "W", "W", "B", "W", "W", "W", "B"},
				{"", "B", "B", "B", "W", "B", "B", "B", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "", "", "B", "W", "B", "", "", ""},
				{"", "", "", "", "B", "", "", "", ""}},
			pt: point.New(4, 4),
			c:  color.Black,
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
			tc.b.FindEnemyGroups(tc.pt, tc.c)
			got := tc.b.String()
			if fmt.Sprintf("%s", got) != fmt.Sprintf("%s", tc.exp) {
				t.Errorf("got %s, but expected %s", got, tc.exp)
			}
		})
	}
}
