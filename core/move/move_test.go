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
		col          color.Color
		exp          *Move
		expErrSubstr string
	}{
		{
			desc:  "Valid Placement",
			sgfPt: "ab",
			col:   color.Black,
			exp:   New(color.Black, point.New(0, 1)),
		},
		{
			desc:  "Pass",
			sgfPt: "",
			col:   color.Black,
			exp:   NewPass(color.Black),
		},
		{
			desc:  "Valid Placement White",
			sgfPt: "ab",
			col:   color.White,
			exp:   New(color.White, point.New(0, 1)),
		},
		{
			desc:  "Pass White",
			sgfPt: "",
			col:   color.White,
			exp:   NewPass(color.White),
		},
		{
			desc:         "Invalid Characters",
			sgfPt:        "!=",
			col:          color.Black,
			expErrSubstr: "only a-zA-Z",
		},
		{
			desc:         "Invalid X",
			sgfPt:        "!b",
			col:          color.Black,
			expErrSubstr: "convert coordinate for x-value",
		},
		{
			desc:         "Invalid Y",
			sgfPt:        "a=",
			col:          color.Black,
			expErrSubstr: "convert coordinate for y-value",
		},
		{
			desc:         "Extra Characters",
			sgfPt:        "abc",
			col:          color.Black,
			expErrSubstr: "two letter char",
		},
		{
			desc:         "One Character",
			sgfPt:        "a",
			col:          color.Black,
			expErrSubstr: "two letter char",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := FromSGFPoint(tc.col, tc.sgfPt)

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
		col          color.Color
		exp          []*Move
		expErrSubstr string
	}{
		{
			desc:      "Valid List",
			sgfPtList: []string{"ab", "cd", "ee"},
			col:       color.Black,
			exp:       []*Move{New(color.Black, point.New(0, 1)), New(color.Black, point.New(2, 3)), New(color.Black, point.New(4, 4))},
		},
		{
			desc:      "Valid List White",
			sgfPtList: []string{"ab", "cd", "ee"},
			col:       color.White,
			exp:       []*Move{New(color.White, point.New(0, 1)), New(color.White, point.New(2, 3)), New(color.White, point.New(4, 4))},
		},
		{
			desc:      "Empty Move List",
			sgfPtList: []string{},
			col:       color.Black,
			exp:       []*Move{},
		},
		{
			desc:         "Contains Pass",
			sgfPtList:    []string{"ab", "", "ef"},
			col:          color.Black,
			expErrSubstr: "non-empty",
		},
		{
			desc:         "Invalid Characters",
			sgfPtList:    []string{"ab", "cd", "%$"},
			col:          color.Black,
			expErrSubstr: "only a-zA-Z",
		},
		{
			desc:         "Invalid X",
			sgfPtList:    []string{"ab", "cd", "(z"},
			col:          color.Black,
			expErrSubstr: "convert coordinate for x-value",
		},
		{
			desc:         "Invalid Y",
			sgfPtList:    []string{"ab", "cd", "g#"},
			col:          color.Black,
			expErrSubstr: "convert coordinate for y-value",
		},
		{
			desc:         "Extra Characters",
			sgfPtList:    []string{"abwegasd", "cd", "hf"},
			col:          color.Black,
			expErrSubstr: "two letter char",
		},
		{
			desc:         "One Character",
			sgfPtList:    []string{"a", "cd", "hf"},
			col:          color.Black,
			expErrSubstr: "two letter char",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := ListFromSGFPoints(tc.col, tc.sgfPtList)

			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			for i, exp := range tc.exp {
				if got[i].color != exp.color || got[i].point.X() != exp.point.X() || got[i].point.Y() != exp.point.Y() {
					t.Errorf("got %v%v, expected %v%v", got[i].color, got[i].point, exp.color, exp.point)
				}
			}
		})
	}
}
