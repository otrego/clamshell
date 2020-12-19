package move

import (
	"testing"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/point"
)

func TestFromSGFPoint(t *testing.T) {
	testCases := []struct {
		desc  string
		sgfPt string
		exp   *Move
	}{
		{
			desc:  "Valid Placement",
			sgfPt: "ab",
			exp:   NewMove(color.Black, point.New(0, 1)),
		},
		{
			desc:  "Pass",
			sgfPt: "",
			exp:   NewPass(color.Black),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			m, _ := FromSGFPoint(color.Black, tc.sgfPt)
			if tc.exp.point != nil {
				if m.color != tc.exp.color || m.point.X() != tc.exp.point.X() || m.point.Y() != tc.exp.point.Y() {
					t.Errorf("got %v%v, expected %v%v", m.color, m.point, tc.exp.color, tc.exp.point)
				}
			} else {
				if m.color != tc.exp.color {
					t.Errorf("got %v, expected %v", m.color, tc.exp.color)
				}
			}

		})
	}
}
