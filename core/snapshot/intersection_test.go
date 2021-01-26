package snapshot

import (
	"testing"

	"github.com/otrego/clamshell/core/snapshot/symbol"
)

func TestIntersection_TopLayerUnicodeChar(t *testing.T) {
	testCases := []struct {
		in  *Intersection
		exp string
	}{
		{
			in: &Intersection{
				Base: symbol.TopLeft,
			},
			exp: "┏",
		},
		{
			in: &Intersection{
				Base:  symbol.TopLeft,
				Stone: symbol.BlackStone,
			},
			exp: "●",
		},
		{
			in: &Intersection{
				Base:  symbol.TopLeft,
				Stone: symbol.BlackStone,
				Mark:  symbol.Xmark,
			},
			exp: "☓",
		},
		{
			in:  &Intersection{},
			exp: " ",
		},
	}

	for _, tc := range testCases {
		if out := tc.in.TopLayerUnicodeString(); out != tc.exp {
			t.Errorf("(%v).TopLayerUnicodeString()=%s, but expected %s", tc.in, out, tc.exp)
		}
	}
}
