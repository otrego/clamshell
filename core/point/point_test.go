package point

import (
	"fmt"
	"testing"

	"github.com/otrego/clamshell/core/errcheck"
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

func TestPointToString(t *testing.T) {
	result := fmt.Sprintf("%v", New(2, 3))
	const expected = "{2,3}"
	if result != expected {
		t.Errorf("Point() to string: expected %v, got %v", expected, result)
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
			in: &Point{
				x: 8,
				y: 16,
			},
			want: "iq",
		},
		{
			desc: "Point => SGF",
			in: &Point{
				x: 12,
				y: 5,
			},
			want: "mf",
		},
		{
			desc: "Point => SGF",
			in: &Point{
				x: 33,
				y: 8,
			},
			want: "Hi",
		},
		{
			desc: "Point => SGF",
			in: &Point{
				x: 40,
				y: 51,
			},
			want: "OZ",
		},
		{
			desc: "Point => SGF, out of range",
			in: &Point{
				x: 100,
				y: 200,
			},
			expErrSubstr: "out of range x or y",
		},
	}

	for _, tc := range testToSGFCases {
		t.Run(tc.desc, func(t *testing.T) {
			toSGFOut, err := New(tc.in.x, tc.in.y).ToSGF()
			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			if toSGFOut != tc.want {
				t.Errorf("%q.ToSGF() = %q, but wanted %q", tc.in,
					toSGFOut, tc.want)
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
			want: &Point{
				x: 8,
				y: 16,
			},
		},
		{
			desc: "SGF => Point",
			in:   "mf",
			want: &Point{
				x: 12,
				y: 5,
			},
		},
		{
			desc: "SGF => Point",
			in:   "Hi",
			want: &Point{
				x: 33,
				y: 8,
			},
		},
		{
			desc: "SGF => Point",
			in:   "OZ",
			want: &Point{
				x: 40,
				y: 51,
			},
		},
		{
			desc:         "SGF => Point: Invalid chars",
			in:           "&&",
			expErrSubstr: "invalid characters",
		},
	}

	for _, tc := range testToPointCases {
		t.Run(tc.desc, func(t *testing.T) {
			toPointOut, err := NewFromSGF(tc.in)
			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Error(cerr)
				return
			}
			if err != nil {
				return
			}

			// include the point.go *Point type X Y getters below
			pointX := toPointOut.X()
			pointY := toPointOut.Y()
			if pointX != tc.want.x {
				t.Errorf("%q.toPointOut.x = %q, but wanted %q", tc.in,
					pointX, tc.want.x)
			}
			if pointY != tc.want.y {
				t.Errorf("%q.toPointOut.y = %q, but wanted %q", tc.in,
					pointY, tc.want.y)
			}
		})
	}

}
