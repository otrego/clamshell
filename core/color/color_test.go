package color

import (
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
	out, err := FromSGFProp("B")
	if err != nil {
		t.Fatal(err)
	}
	if out != Black {
		t.Errorf("got %v, but wanted %v", out, Black)
	}

	out, err = FromSGFProp("AB")
	if err != nil {
		t.Fatal(err)
	}
	if out != Black {
		t.Errorf("got %v, but wanted %v", out, Black)
	}

	out, err = FromSGFProp("W")
	if err != nil {
		t.Fatal(err)
	}
	if out != Black {
		t.Errorf("got %v, but wanted %v", out, White)
	}

	out, err = FromSGFProp("AW")
	if err != nil {
		t.Fatal(err)
	}
	if out != Black {
		t.Errorf("got %v, but wanted %v", out, White)
	}
}
