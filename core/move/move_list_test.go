package move

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/point"
)

func TestMoveList_String(t *testing.T) {
	mvStr := List{
		New(color.Black, point.New(1, 2)),
		New(color.White, point.New(3, 4)),
	}.String()

	exp := "{{B, {1,2}}, {W, {3,4}}}"
	if mvStr != exp {
		t.Errorf("mvlist.String()=%s, but expected %s. diff=%s", mvStr, exp, cmp.Diff(mvStr, exp))
	}
}

func TestMoveList_Sort(t *testing.T) {
	mv := List{
		New(color.Black, point.New(1, 2)),
		New(color.White, point.New(3, 4)),
		New(color.Black, point.New(2, 2)),
		New(color.White, point.New(2, 2)),
	}
	exp := List{
		New(color.Black, point.New(1, 2)),
		New(color.Black, point.New(2, 2)),
		New(color.White, point.New(2, 2)),
		New(color.White, point.New(3, 4)),
	}

	mv.Sort()
	if !reflect.DeepEqual(mv, exp) {
		t.Errorf("after sorting got %v, but expected %v. diff=%v", mv, exp,
			cmp.Diff(mv.String(), exp.String()))
	}
}
