package move

import (
	"testing"

	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/errcheck"
	"github.com/otrego/clamshell/core/point"
)

func TestFromSGFPoint(t *testing.T) {
	testCases := []struct {
		desc         string
		sgfPt        string
		exp          *Move
		expErrSubstr string
	}{
		{
			desc:  "Valid Placement",
			sgfPt: "ab",
			exp:   NewMove(color.Black, point.New(0, 1)),
		},
		{
			desc:  "Valid Placement",
			sgfPt: "!=",
			exp:   NewMove(color.Black, point.New(5, 1)),
		},
		{
			desc:  "Pass",
			sgfPt: "",
			exp:   NewPass(color.Black),
		},
		{
			desc:         "Out of bounds",
			sgfPt:        "--",
			expErrSubstr: "SGF string x and y value entries must",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := FromSGFPoint(color.Black, tc.sgfPt)

			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			if tc.exp.point != nil {
				if got.color != tc.exp.color || got.point.X() != tc.exp.point.X() || got.point.Y() != tc.exp.point.Y() {
					t.Errorf("got %v%v, expected %v%v", got.color, got.point, tc.exp.color, tc.exp.point)
				}
			} else {
				if got.color != tc.exp.color {
					t.Errorf("got %v, expected %v", got.color, tc.exp.color)
				}
			}

		})
	}
}

func TestListFromSGFPoints(t *testing.T) {
	testCases := []struct {
		desc         string
		sgfPt        string
		exp          *Move
		expErrSubstr string
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
		{
			desc:         "Out of bounds",
			sgfPt:        "--",
			expErrSubstr: "SGF string x and y value entries must",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := FromSGFPoint(color.Black, tc.sgfPt)

			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			if tc.exp.point != nil {
				if got.color != tc.exp.color || got.point.X() != tc.exp.point.X() || got.point.Y() != tc.exp.point.Y() {
					t.Errorf("got %v%v, expected %v%v", got.color, got.point, tc.exp.color, tc.exp.point)
				}
			} else {
				if got.color != tc.exp.color {
					t.Errorf("got %v, expected %v", got.color, tc.exp.color)
				}
			}

		})
	}
}
