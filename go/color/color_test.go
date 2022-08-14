package color

import (
	"errors"
	"testing"
)

func TestOpposite(t *testing.T) {
	testCases := []struct {
		desc string
		in   string
		want string
	}{
		{
			desc: "white=>black",
			in:   "W",
			want: "B",
		},
		{
			desc: "black=>white",
			in:   "B",
			want: "W",
		},
		{
			desc: "empty=>empty",
			in:   "",
			want: "",
		},
		{
			desc: "any=>any",
			in:   "any",
			want: "any",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			out := Color(tc.in).Opposite()
			if out != Color(tc.want) {
				t.Errorf("%q.Opposite()=%q, but wanted %q", tc.in, out, tc.want)
			}
		})
	}
}

func TestFromSGFProp(t *testing.T) {
	testCases := []struct {
		desc       string
		in         string
		want       Color
		expErrType error
	}{
		{
			desc:       "B=>Black",
			in:         "B",
			want:       Black,
			expErrType: nil,
		},
		{
			desc:       "AB=>Black",
			in:         "AB",
			want:       Black,
			expErrType: nil,
		},
		{
			desc:       "W=>White",
			in:         "W",
			want:       White,
			expErrType: nil,
		},
		{
			desc:       "AW=>White",
			in:         "AW",
			want:       White,
			expErrType: nil,
		},
		{
			desc:       "empty=>empty",
			in:         "",
			want:       Empty,
			expErrType: ErrColorConversion,
		},
		{
			desc:       "any=>any",
			in:         "any",
			want:       Empty,
			expErrType: ErrColorConversion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			out, err := FromSGFProp(tc.in)
			if out != tc.want {
				t.Errorf("%v was passed, got %v, wanted %v", tc.in, out, tc.want)
			}
			if !errors.Is(err, tc.expErrType) {
				t.Errorf("Got err %v, but  expected error of type %v", err, tc.expErrType)
			}
		})
	}
}
