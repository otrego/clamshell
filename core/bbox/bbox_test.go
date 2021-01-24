package bbox

import (
	"testing"

	"github.com/otrego/clamshell/core/point"
)

func TestBboxBasics(t *testing.T) {
	tl := point.New(1, 2)
	br := point.New(10, 13)

	bbox, err := New(tl, br)
	if err != nil {
		t.Fatalf("error making bounding box: %v", err)
	}

	var exp int
	exp = 1
	if got := bbox.TopLeft().X(); got != exp {
		t.Errorf("bbox.TopLeft().X()=%d, but expected %d", got, exp)
	}
	if got := bbox.Left(); got != exp {
		t.Errorf("bbox.Left()=%d, but expected %d", got, exp)
	}

	exp = 2
	if got := bbox.TopLeft().Y(); got != exp {
		t.Errorf("bbox.TopLeft().Y()=%d, but expected %d", got, exp)
	}
	if got := bbox.Top(); got != exp {
		t.Errorf("bbox.Top()=%d, but expected %d", got, exp)
	}

	exp = 10
	if got := bbox.BotRight().X(); got != exp {
		t.Errorf("bbox.BotRight().X()=%d, but expected %d", got, exp)
	}
	if got := bbox.Right(); got != exp {
		t.Errorf("bbox.Right()=%d, but expected %d", got, exp)
	}

	exp = 13
	if got := bbox.BotRight().Y(); got != exp {
		t.Errorf("bbox.BotRight().Y()=%d, but expected %d", got, exp)
	}
	if got := bbox.Bottom(); got != exp {
		t.Errorf("bbox.Bottom()=%d, but expected %d", got, exp)
	}

	exp = 9
	if got := bbox.Width(); got != exp {
		t.Errorf("bbox.Width()=%d, but expected %d", got, exp)
	}

	exp = 11
	if got := bbox.Height(); got != exp {
		t.Errorf("bbox.Height()=%d, but expected %d", got, exp)
	}
}
