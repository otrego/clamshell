package point

import (
	"fmt"
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

// The below function and related elements were contributed by
// un nuovo amico GitHub dall'Italia, di nome ilmanzo
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
		desc          string
		in            *Point
		want          string
		wantErrSubstr string
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
		// // Error Cases
		// {
		// 	desc:          "Error Point => SGF",
		// 	in:            New(40, 52),
		// 	wantErrSubstr: "Point Y value > 51, out of range",
		// },
	}

	for _, tc := range testToSGFCases {
		t.Run(tc.desc, func(t *testing.T) {
			sgfOut, _ := New(tc.in.x, tc.in.y).ToSGF()
			if sgfOut != tc.want {
				t.Errorf("%q.ToSGF() = %q, but wanted %q", tc.in,
					sgfOut, tc.want)
			}
		})
	}
}

func TestSGFToPointTranslate(t *testing.T) {
	testToPointCases := []struct {
		desc          string
		in            string
		want          *Point
		wantErrSubstr string
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
	}

	for _, tc := range testToPointCases {
		t.Run(tc.desc, func(t *testing.T) {
			pntOut, _ := NewFromSGF(tc.in)
			// Utilizing the point.go *Point type X Y getters below
			x := pntOut.X()
			y := pntOut.Y()
			if x != tc.want.x {
				t.Errorf("%q.pntOut.x = %q, but wanted %q", tc.in,
					x, tc.want.x)
			}
			if y != tc.want.y {
				t.Errorf("%q.pntOut.y = %q, but wanted %q", tc.in,
					y, tc.want.y)
			}
		})
	}
}
