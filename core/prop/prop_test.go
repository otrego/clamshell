package prop_test

import (
	"github.com/otrego/clamshell/core/prop"
	"testing"
)

func TestValidate(t *testing.T) {
	valid := prop.Field("TB")

	test := prop.Validate(valid)
	if test != true {
		t.Errorf("Validate(%v) SGF field returned %v; want true", valid, test)
	}
}

func TestInvalidate(t *testing.T) {
	invalid := prop.Field("TQB")

	test := prop.Validate(invalid)
	if test != false {
		t.Errorf("Validate(%v) SGF field returned %v; want false", invalid, test)
	}
}
