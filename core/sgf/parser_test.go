package sgf_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/otrego/clamshell/core/errcheck"
	"github.com/otrego/clamshell/core/sgf"
	"github.com/otrego/clamshell/core/treepath"
)

type propmap map[string][]string

func TestParse(t *testing.T) {
	testCases := []struct {
		desc         string
		sgf          string
		pathToProps  map[string]propmap
		expErrSubstr string
	}{
		{
			desc: "basic parse, no errors",
			sgf:  "(;GM[1])",
		},
		{
			desc: "check parsing root: simple property",
			sgf:  "(;GM[1])",
			pathToProps: map[string]propmap{
				".": propmap{
					"GM": []string{"1"},
				},
			},
		},
		{
			desc: "check parsing root: lots of whitespace",
			sgf:  "\n   (;  C[1 \n1]   )",
			pathToProps: map[string]propmap{
				".": propmap{
					"C": []string{"1 \n1"},
				},
			},
		},
		{
			desc: "check parsing root: multi property",
			sgf:  "(;GM[1]AW[ab][bc][dc])",
			pathToProps: map[string]propmap{
				".": propmap{
					"GM": []string{"1"},
					"AW": []string{"ab", "bc", "dc"},
				},
			},
		},
		{
			desc: "add new node",
			sgf:  "(;GM[1]AW[ab][bc][dc];B[cc])",
			pathToProps: map[string]propmap{
				".": propmap{
					"GM": []string{"1"},
					"AW": []string{"ab", "bc", "dc"},
				},
				"0": propmap{
					"B": []string{"cc"},
				},
			},
		},
		{
			desc: "comment with escaped rbrace",
			sgf:  `(;C[aoeu [1k\]])`,
			pathToProps: map[string]propmap{
				".": propmap{
					"C": []string{"aoeu [1k]"},
				},
			},
		},
		{
			desc: "comment with non-escapng backslash",
			sgf:  `(;C[aoeu \z])`,
			pathToProps: map[string]propmap{
				".": propmap{
					"C": []string{"aoeu \\z"},
				},
			},
		},
		{
			desc: "comment with lots of escaping",
			sgf:  `(;C[\\\]])`,
			pathToProps: map[string]propmap{
				".": propmap{
					"C": []string{`\\]`},
				},
			},
		},
		{
			desc: "basic variation",
			sgf:  `(;GM[1](;B[aa]W[ab])(;B[ab]W[ac]))`,
			pathToProps: map[string]propmap{
				"0.0": propmap{
					"W": []string{"ab"},
				},
				"1.0": propmap{
					"W": []string{"ac"},
				},
			},
		},

		{
			desc: "mega mark test",
			sgf: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[0.00]
PW[White]PB[Black]
AW[na][oa][pa][qa][ra][sa][ka][la][ma][ja]
AB[nb][ob][pb][qb][rb][sb][kb][lb][mb][jb]
LB[pa:A][ob:2][pb:B][pc:C][pd:D]
[oa:1][oc:3][ne:9][oe:8][pe:7][qe:6][re:5][se:4]
[nf:15][of:14][pf:13][qf:11][rf:12][sf:10]
[ng:22][og:44][pg:100]
[ka:a][kb:b][kc:c][kd:d][ke:e][kf:f][kg:g]
[ma:\u4e00][mb:\u4e8c][mc:\u4e09][md:\u56db][me:\u4e94]
[la:\u516d][lb:\u4e03][lc:\u516b][ld:\u4e5d][le:\u5341]
MA[na][nb][nc]
CR[qa][qb][qc]
TR[sa][sb][sc]
SQ[ra][rb][rc]
)`,
			pathToProps: map[string]propmap{
				".": propmap{
					"GM": []string{"1"},
					"SQ": []string{"ra", "rb", "rc"},
					"PW": []string{"White"},
					"SZ": []string{"19"},
				},
			},
		},

		{
			desc: "complex problem",
			sgf: `
(;GM[1]FF[4]CA[UTF-8]AP[Glift]ST[2]
RU[Japanese]SZ[19]KM[0.00]
C[Black to play. There aren't many options
to choose from, but you might be surprised at the answer!]
PW[White]PB[Black]AW[pa][qa][nb][ob][qb][oc][pc][md][pd][ne][oe]
AB[na][ra][mb][rb][lc][qc][ld][od][qd][le][pe][qe][mf][nf][of][pg]
(;B[mc]
	;W[nc]C[White lives.])
(;B[ma]
	(;W[oa]
		;B[nc]
		;W[nd]
		;B[mc]C[White dies.]GB[1])
	(;W[mc]
		(;B[oa]
		;W[nd]
		;B[pb]C[White lives])
		(;B[nd]
			;W[nc]
			;B[oa]C[White dies.]GB[1]))
	(;W[nd]
		;B[mc]
		;W[oa]
		;B[nc]C[White dies.]GB[1]))
(;B[nc]
	;W[mc]C[White lives])
(;B[]C[A default consideration]
	;W[mc]C[White lives easily]))`,
			pathToProps: map[string]propmap{
				".": propmap{
					"GM": []string{"1"},
				},
				"0.0": propmap{
					"W": []string{"nc"},
				},
				// should be same, since treepath terminates
				"0.0.0": propmap{
					"W": []string{"nc"},
				},
				"1.1.0": propmap{
					"B": []string{"oa"},
				},
				"3": propmap{
					"B": []string{""},
				},
			},
		},

		// error cases
		{
			desc:         "error parsing: nonsense",
			sgf:          "I am a banana",
			expErrSubstr: "unexpected char",
		},
		{
			desc:         "error parsing: unclosed",
			sgf:          "(;C",
			expErrSubstr: "expected to end on ')'",
		},
		{
			desc:         "error parsing: bad property",
			sgf:          "(;C)",
			expErrSubstr: "during property parsing",
		},
		{
			desc:         "error parsing: empty variation",
			sgf:          "(;C[])())",
			expErrSubstr: "empty variation",
		},
		{
			desc:         "error parsing: unclosed property data",
			sgf:          "(;C[)())",
			expErrSubstr: "ended in nested condition",
		},
		{
			desc:         "error parsing: bad property",
			sgf:          "(;weird[])",
			expErrSubstr: "unexpected character between",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := sgf.FromString(tc.sgf).Parse()
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

			for path, pmap := range tc.pathToProps {
				tp, err := treepath.Parse(path)
				if err != nil {
					t.Error(err)
					continue
				}
				n := tp.Apply(g.Root)
				for prop, expData := range pmap {
					foundData, ok := n.Properties[prop]
					if !ok {
						t.Errorf("At path %q, properties did not contain expected property key %q. Properties were %v", path, prop, n.Properties)
					}
					if !cmp.Equal(foundData, expData) {
						t.Errorf("At path %q, property %q was %v, but expected %v", path, prop, foundData, expData)
					}
				}
			}
		})
	}
}
