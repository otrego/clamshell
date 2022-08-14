package point

import (
	"github.com/google/go-cmp/cmp"

	"testing"
)

func TestCoordEqual(t *testing.T) {
	pt := NewCoord(1.2, 3.4)
	other := &Coord{
		x: 1.2,
		y: 3.4,
	}

	if pt == other {
		t.Errorf("expected points to have different references")
	}
	if !pt.Equal(other) {
		t.Errorf("expected pt %v to be equal to %v", pt, other)
	}
	if !cmp.Equal(pt, other) {
		t.Errorf("expected pt %v to be equal to %v using the cmp library", pt, other)
	}
}
