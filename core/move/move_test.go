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
			desc:  "Pass",
			sgfPt: "",
			exp:   NewPass(color.Black),
		},
		{
			desc:         "Invalid Characters",
			sgfPt:        "!=",
			expErrSubstr: "only a-zA-Z",
		},
		{
			desc:         "Invalid X",
			sgfPt:        "!b",
			expErrSubstr: "convert coordinate for x-value",
		},
		{
			desc:         "Invalid Y",
			sgfPt:        "a=",
			expErrSubstr: "convert coordinate for y-value",
		},
		{
			desc:         "Extra Characters",
			sgfPt:        "abc",
			expErrSubstr: "two letter char",
		},
		{
			desc:         "One Character",
			sgfPt:        "a",
			expErrSubstr: "two letter char",
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
		sgfPtList    []string
		exp          []*Move
		expErrSubstr string
	}{
		{
			desc:      "Valid List",
			sgfPtList: []string{"ab", "cd", "ee"},
			exp:       []*Move{NewMove(color.Black, point.New(0, 1)), NewMove(color.Black, point.New(2, 3)), NewMove(color.Black, point.New(4, 4))},
		},
		{
			desc:      "Empty Move List",
			sgfPtList: []string{},
			exp:       []*Move{},
		},
		{
			desc:         "Contains Pass",
			sgfPtList:    []string{"ab", "", "ef"},
			expErrSubstr: "non-empty",
		},
		{
			desc:         "Invalid Characters",
			sgfPtList:    []string{"ab", "cd", "%$"},
			expErrSubstr: "only a-zA-Z",
		},
		{
			desc:         "Invalid X",
			sgfPtList:    []string{"ab", "cd", "(z"},
			expErrSubstr: "convert coordinate for x-value",
		},
		{
			desc:         "Invalid Y",
			sgfPtList:    []string{"ab", "cd", "g#"},
			expErrSubstr: "convert coordinate for y-value",
		},
		{
			desc:         "Extra Characters",
			sgfPtList:    []string{"abwegasd", "cd", "hf"},
			expErrSubstr: "two letter char",
		},
		{
			desc:         "One Character",
			sgfPtList:    []string{"a", "cd", "hf"},
			expErrSubstr: "two letter char",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := ListFromSGFPoints(color.Black, tc.sgfPtList)

			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			for i, exp := range tc.exp {
				if exp.point != nil {
					if got[i].color != exp.color || got[i].point.X() != exp.point.X() || got[i].point.Y() != exp.point.Y() {
						t.Errorf("got %v%v, expected %v%v", got[i].color, got[i].point, exp.color, exp.point)
					}
				} else {
					if got[i].color != exp.color {
						t.Errorf("got %v, expected %v", got[i].color, exp.color)
					}
				}
			}
		})
	}
}
