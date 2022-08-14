package prop_test

import (
	"testing"

	"github.com/otrego/clamshell/go/prop"
)

func TestValidate(t *testing.T) {
	valid := prop.Prop("TB")

	test := prop.Validate(valid)
	if test != true {
		t.Errorf("Validate(%v) returned %v; want true", valid, test)
	}
}

func TestInvalidate(t *testing.T) {
	invalid := prop.Prop("TQB")

	test := prop.Validate(invalid)
	if test != false {
		t.Errorf("Validate(%v) returned %v; want false", invalid, test)
	}
}
