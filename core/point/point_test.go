package point

import "testing"

func TestCreate(t *testing.T) {
	p := New(1, 2)

	if gotX := p.X(); gotX != 1 {
		t.Errorf("p.X()=%v, expected %v", gotX, 1)
	}

	if gotY := p.Y(); gotY != 2 {
		t.Errorf("p.Y()=%v, expected %v", gotY, 2)
	}
}

func TestPointSGFTranslation(t *testing.T) {
	type Point struct {
		x int64
		y int64
	}

	// First test translation from integer-point to SGF-string-point
	testToSGFCases := []struct {
		desc string
		in   Point
		want string
	}{
		{
			desc: "Point => SGF",
			in: Point{
				x: 8,
				y: 16,
			},
			want: "iq",
		},
		{
			desc: "Point => SGF",
			in: Point{
				x: 12,
				y: 5,
			},
			want: "mf",
		},
		{
			desc: "Point => SGF",
			in: Point{
				x: 33,
				y: 8,
			},
			want: "Hi",
		},
		{
			desc: "Point => SGF",
			in: Point{
				x: 40,
				y: 51,
			},
			want: "OZ",
		},
	}

	for _, tc := range testToSGFCases {
		t.Run(tc.desc, func(t01 *testing.T) {
			toSGFOut := New(tc.in.x, tc.in.y).ToSGF()
			if t01 == nil || toSGFOut != tc.want {
				t01.Errorf("%q.ToSGF() = %q, but wanted %q", tc.in,
					toSGFOut, tc.want)
			}
		})
	}

	// Second test translation from SGF-string-point to integer-point
	testToPointCases := []struct {
		desc string
		in   string
		want *Point
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
	}

	for _, tc := range testToPointCases {
		t.Run(tc.desc, func(t01 *testing.T) {
			toPointOut := NewFromSGF(tc.in)
			// include the point.go *Point type X Y getters below
			pointX := toPointOut.X()
			pointY := toPointOut.Y()
			if t01 == nil || pointX != tc.want.x {
				t01.Errorf("%q.toPointOut.x = %q, but wanted %q", tc.in,
					pointX, tc.want.x)
			}
			if t01 == nil || pointY != tc.want.y {
				t01.Errorf("%q.toPointOut.y = %q, but wanted %q", tc.in,
					pointY, tc.want.y)
			}
		})
	}
}
