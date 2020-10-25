package game

import (
	"testing"

	"github.com/otrego/clamshell/core/color"
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
