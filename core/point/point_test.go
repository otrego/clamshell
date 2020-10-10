package point

import (
	"fmt"
	"github.com/otrego/clamshell/core/errcheck"
	"testing"
)

func TestCreate(t *testing.T) {
	p := New(1, 2)

	if gotX := p.X(); gotX != 1 {
		t.Errorf("p.X()=%v, expected %v", gotX, 1)
	}

	if gotY := p.Y(); gotY != 2 {
		t.Errorf("p.Y()=%v, expected %v", gotY, 2)
	}
}

// The below function TestPointToString() was authored by ilmanzo
func TestPointToString(t *testing.T) {
	result := fmt.Sprintf("%v", New(2, 3))
	const expected = "{2,3}"
	if result != expected {
		t.Errorf("Point() to string: expected %v, got %v", expected,
			result)
	}
}

func TestPointToSGFTranslate(t *testing.T) {
	// First test translation from integer-point to SGF-string-point
	testToSGFCases := []struct {
		desc         string
		in           *Point
		want         string
		expErrSubstr string
	}{
		{
			desc: "Point => SGF",
			in:   New(8, 16),
			want: "iq",
		},
		{
			desc: "Point => SGF",
			in:   New(12, 5),
			want: "mf",
		},
		{
			desc: "Point => SGF",
			in:   New(33, 8),
			want: "Hi",
		},
		{
			desc: "Point => SGF",
			in:   New(40, 51),
			want: "OZ",
		},
		// Error Cases:
		// Negative X Y values
		{
			desc:         "Error: Point => SGF: X value < 0",
			in:           New(-33, 8),
			expErrSubstr: "Point X value < 0, out of range",
		},
		{
			desc:         "Error: Point => SGF: Y value < 0",
			in:           New(33, -8),
			expErrSubstr: "Point Y value < 0, out of range",
		},
		// Out of range X Y values
		{
			desc:         "Error: Point => SGF: X value > 51",
			in:           New(53, 8),
			expErrSubstr: "Point X value > 51, out of range",
		},
		{
			desc:         "Error: Point => SGF: Y value > 51",
			in:           New(33, 53),
			expErrSubstr: "Point Y value > 51, out of range",
		},
	}

	// Test cases
	for _, tc := range testToSGFCases {
		t.Run(tc.desc, func(t *testing.T) {
			sgfOut, err := New(tc.in.x, tc.in.y).ToSGF()
			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if sgfOut != tc.want {
				t.Errorf("%q.ToSGF() = %q, but wanted %q", tc.in,
					sgfOut, tc.want)
				// if cerr != nil && sgfOut != tc.want {
				if cerr != nil {
					t.Fatal(cerr)
				}
				if err != nil {
					return
				}
				if sgfOut == "" {
					t.Fatal("unexpected empty SGF string value")
				}
			}
		})
	}
}

func TestSGFToPointTranslate(t *testing.T) {
	testToPointCases := []struct {
		desc         string
		in           string
		want         *Point
		expErrSubstr string
	}{
		{
			desc: "SGF => Point",
			in:   "iq",
			want: New(8, 16),
		},
		{
			desc: "SGF => Point",
			in:   "mf",
			want: New(12, 5),
		},
		{
			desc: "SGF => Point",
			in:   "Hi",
			want: New(33, 8),
		},
		{
			desc: "SGF => Point",
			in:   "OZ",
			want: New(40, 51),
		},
		// Error Cases:
		// Empty string
		{
			desc:         "Error: SGF => Point empty string",
			in:           "",
			expErrSubstr: "Error: SGF X Y string values empty",
		},
		// Short string
		{
			desc:         "Error: SGF => Point short string",
			in:           "Q",
			expErrSubstr: "SGF X or Y string value missing",
		},
		// Empty string
		{
			desc:         "Error: SGF => Point long string",
			in:           "xyZ",
			expErrSubstr: "SGF X Y string values empty",
		},
	}

	for _, tc := range testToPointCases {
		t.Run(tc.desc, func(t *testing.T) {
			pntOut, err := NewFromSGF(tc.in)
			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			// Utilizing the point.go *Point type X Y getters below
			if pntOut != tc.want {
				if pntOut != nil {
					x := pntOut.X()
					y := pntOut.Y()
					if x != tc.want.x {
						t.Errorf("%q.pntOut.x = %q, but wanted %q",
							tc.in, x, tc.want.x)
					}
					if y != tc.want.y {
						t.Errorf("%q.pntOut.y = %q, but wanted %q",
							tc.in,
							y, tc.want.y)
					}
					if cerr != nil {
						t.Fatal(cerr)
					}
					if err != nil {
						return
					}
				}
				if pntOut == nil {
					t.Fatal("unexpected nil point integer value")
				}
			}
		})
	}
}
