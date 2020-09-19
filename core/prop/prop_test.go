package prop_test

import (
	"github.com/otrego/clamshell/core/prop"
	"testing"
)

func TestValidate(t *testing.T) {

	var valid prop.Field = "TB"

	test := prop.Validate(valid)
	if test != true {
		t.Errorf("Valid SGF field returned %v; want true", test)
	}
}

func TestInvalidate(t *testing.T) {

	var invalid prop.Field = "TQB"

	test := prop.Validate(invalid)
	if test != false {
		t.Errorf("Invalid SGF field returned %v; want false", test)
	}
}
