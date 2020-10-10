// Package errcheck checks error cases.
package errcheck

import (
	"fmt"
	"strings"
)

// CheckCases checks error cases for tests
func CheckCases(err error, expErrSubstr string) error {
	if err == nil && expErrSubstr != "" {
		return fmt.Errorf("got no error but expected error containing %q",
			expErrSubstr)
	} else if err != nil && expErrSubstr == "" {
		return fmt.Errorf("got error %q but expected no error",
			err.Error())
	} else if err != nil && !strings.Contains(err.Error(),
		expErrSubstr) {
		return fmt.Errorf("Got error %q but expected it to contain %q",
			err.Error(), expErrSubstr)
	}
	return nil
}
