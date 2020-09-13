package sgf

import (
	"testing"

	"github.com/otrego/clamshell/core/errcheck"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		desc         string
		sgf          string
		expErrSubstr string
	}{
		{
			desc: "basic parse, no errors",
			sgf:  "(;GM[1])",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := FromString(tc.sgf).Parse()
			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Fatal(cerr)
			}
			if err != nil {
				return
			}
			if g == nil {
				t.Fatal("unexpectedly nil game")
			}
		})
	}
}
