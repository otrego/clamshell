package movetree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefaults(t *testing.T) {
	g := New()

	expSize := 19
	if got := g.Root.GameInfo.Size; got != expSize {
		t.Errorf("g.Root.GameInfo.Size was %v; expected %v", got, expSize)
	}

	expGM := []string{"1"}
	if got := g.Root.SGFProperties["GM"]; !cmp.Equal(got, expGM) {
		t.Errorf("g.Root.SGFProperties[GM] was %v; expected %v", got, expGM)
	}

	expFF := []string{"4"}
	if got := g.Root.SGFProperties["FF"]; !cmp.Equal(got, expFF) {
		t.Errorf("g.Root.SGFProperties[FF] was %v; expected %v", got, expFF)
	}

	expCA := []string{"UTF-8"}
	if got := g.Root.SGFProperties["CA"]; !cmp.Equal(got,expCA) {
		t.Errorf("g.Root.SGFProperties[CA] was %v; expected %v", got, expCA)
	}
}
