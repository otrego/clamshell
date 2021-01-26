package symbol

import "testing"

func TestUnicodeString(t *testing.T) {
	if got := TopLeft.UnicodeString(); got != "┏" {
		t.Errorf("expected TopLeft.UnicodeString() == \"┏\", but was %q", got)
	}

	if got := Symbol(100).UnicodeString(); got != "?" {
		t.Errorf("expected unknown.UnicodeString() == \"?\", but was %q", got)
	}
}
