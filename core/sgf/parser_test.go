package sgf_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/otrego/clamshell/core/color"
	"github.com/otrego/clamshell/core/errcheck"
	"github.com/otrego/clamshell/core/move"
	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/point"
	"github.com/otrego/clamshell/core/sgf"
)

type propmap map[string][]string

type nodeCheck func(n *movetree.Node) error

func TestParse(t *testing.T) {
	testCases := []struct {
		desc            string
		sgf             string
		pathToProps     map[string]propmap
		pathToNodeCheck map[string]nodeCheck
		expErrSubstr    string
	}{
		{
			desc: "basic parse, no errors",
			sgf:  "(;GM[1])",
		},
		{
			desc: "check parsing root: simple property",
			sgf:  "(;GM[1])",
			pathToProps: map[string]propmap{
				"-": propmap{
					"GM": []string{"1"},
				},
			},
		},
		{
			desc: "check parsing root: lots of whitespace",
			sgf:  "\n   (;  ZZ[1 \n1]   )",
			pathToProps: map[string]propmap{
				"-": propmap{
					"ZZ": []string{"1 \n1"},
				},
			},
		},
		{
			desc: "check parsing root: multi property",
			sgf:  "(;GM[1]AW[ab][bc])",
			pathToProps: map[string]propmap{
				"-": propmap{
					"GM": []string{"1"},
				},
			},
			pathToNodeCheck: map[string]nodeCheck{
				"-": func(n *movetree.Node) error {
					expPlacements := []*move.Move{
						move.New(color.White, point.New(0, 1)),
						move.New(color.White, point.New(1, 2)),
					}
					if !reflect.DeepEqual(n.Placements, expPlacements) {
						return fmt.Errorf("incorrect placements; got %v, but wanted %v", n.Placements, expPlacements)
					}
					return nil
				},
			},
		},
		{
			desc: "add new node",
			sgf:  "(;GM[1];B[cc])",
			pathToProps: map[string]propmap{
				"-": propmap{
					"GM": []string{"1"},
				},
			},
			pathToNodeCheck: map[string]nodeCheck{
				"0": func(n *movetree.Node) error {
					expMove := move.New(color.Black, point.New(2, 2))
					if !reflect.DeepEqual(n.Move, expMove) {
						return fmt.Errorf("incorrect move; got %v, but wanted %v", n.Move, expMove)
					}
					return nil
				},
			},
		},
		{
			desc: "comment with escaped rbrace",
			sgf:  `(;ZZ[aoeu [1k\]])`,
			pathToProps: map[string]propmap{
				"-": propmap{
					"ZZ": []string{"aoeu [1k]"},
				},
			},
		},
		{
			desc: "comment with non-escapng backslash",
			sgf:  `(;ZZ[aoeu \z])`,
			pathToProps: map[string]propmap{
				"-": propmap{
					"ZZ": []string{"aoeu \\z"},
				},
			},
		},
		{
			desc: "comment with lots of escaping",
			sgf:  `(;ZZ[\\\]])`,
			pathToProps: map[string]propmap{
				"-": propmap{
					"ZZ": []string{`\\]`},
				},
			},
		},
		{
			desc: "basic variation",
			sgf:  `(;GM[1](;B[aa];W[ab])(;B[ab];W[ac]))`,
			pathToNodeCheck: map[string]nodeCheck{
				"0-0": func(n *movetree.Node) error {
					expMove := move.New(color.White, point.New(0, 1))
					if !reflect.DeepEqual(n.Move, expMove) {
						return fmt.Errorf("incorrect move; got %v, but wanted %v", n.Move, expMove)
					}
					return nil
				},
				"1-0": func(n *movetree.Node) error {
					expMove := move.New(color.White, point.New(0, 2))
					if !reflect.DeepEqual(n.Move, expMove) {
						return fmt.Errorf("incorrect move; got %v, but wanted %v", n.Move, expMove)
					}
					return nil
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
				"-": propmap{
					"GM": []string{"1"},
					"SQ": []string{"ra", "rb", "rc"},
					"PW": []string{"White"},
				},
			},
			pathToNodeCheck: map[string]nodeCheck{
				"-": func(n *movetree.Node) error {
					expSize := 19
					if n.GameInfo.Size != expSize {
						return fmt.Errorf("incorrect size; got %v, but wanted %v", n.GameInfo.Size, expSize)
					}
					return nil
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
				"-": propmap{
					"GM": []string{"1"},
				},
			},
			pathToNodeCheck: map[string]nodeCheck{
				"0-0": func(n *movetree.Node) error {
					expMove := move.New(color.White, point.New(13, 2))
					if !reflect.DeepEqual(n.Move, expMove) {
						return fmt.Errorf("incorrect move; got %v, but wanted %v", n.Move, expMove)
					}
					return nil
				},
				// should be same, since treepath terminates
				"0-0-0": func(n *movetree.Node) error {
					expMove := move.New(color.White, point.New(13, 2))
					if !reflect.DeepEqual(n.Move, expMove) {
						return fmt.Errorf("incorrect move; got %v, but wanted %v", n.Move, expMove)
					}
					return nil
				},
				"1-1-0": func(n *movetree.Node) error {
					expMove := move.New(color.Black, point.New(14, 0))
					if !reflect.DeepEqual(n.Move, expMove) {
						return fmt.Errorf("incorrect move; got %v, but wanted %v", n.Move, expMove)
					}
					return nil
				},
				"3": func(n *movetree.Node) error {
					expMove := move.NewPass(color.Black)
					if !reflect.DeepEqual(n.Move, expMove) {
						return fmt.Errorf("incorrect move; got %v, but wanted %v", n.Move, expMove)
					}
					return nil
				},
			},
		},

		// error cases
		{
			desc:         "basic variation error: two moves on one node",
			sgf:          `(;GM[1](;B[aa]W[ab])(;B[ab];W[ac]))`,
			expErrSubstr: "found two moves",
		},
		{
			desc:         "error parsing: nonsense",
			sgf:          "I am a banana",
			expErrSubstr: "unexpected char",
		},
		{
			desc:         "error parsing: unclosed",
			sgf:          "(;C",
			expErrSubstr: "expected to end on root branch",
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
				t.Fatal("unexpectedly nil movetree")
			}

			for path, pmap := range tc.pathToProps {
				tp, err := movetree.ParsePath(path)
				if err != nil {
					t.Error(err)
					continue
				}
				n := tp.Apply(g.Root)
				for prop, expData := range pmap {
					foundData, ok := n.SGFProperties[prop]
					if !ok {
						t.Errorf("At path %q, properties did not contain expected property key %q. Properties were %v", path, prop, n.SGFProperties)
					}
					if !cmp.Equal(foundData, expData) {
						t.Errorf("At path %q, property %q was %v, but expected %v", path, prop, foundData, expData)
					}
				}
			}
			for path, check := range tc.pathToNodeCheck {
				tp, err := movetree.ParsePath(path)
				if err != nil {
					t.Error(err)
					continue
				}
				n := tp.Apply(g.Root)
				if err := check(n); err != nil {
					t.Errorf("At path %q, found incorrect contents %v", path, err)
				}
			}
		})
	}
}

type propGetter func(*movetree.Node) interface{}

func TestPropertyPostProcessing(t *testing.T) {
	testCases := []struct {
		desc   string
		sgf    string
		path   string
		getter propGetter
		want   interface{}
	}{
		{
			desc: "black move",
			sgf:  "(;GM[1];B[ab])",
			path: "0",
			getter: func(n *movetree.Node) interface{} {
				return n.Move
			},
			want: move.New(color.Black, point.New(0, 1)),
		},
		{
			desc: "white move",
			sgf:  "(;GM[1];W[ab])",
			path: "0",
			getter: func(n *movetree.Node) interface{} {
				return n.Move
			},
			want: move.New(color.White, point.New(0, 1)),
		},
		{
			desc: "black & white placements",
			sgf:  "(;GM[1];AB[ab][ac]AW[bb][bc])",
			path: "0",
			getter: func(n *movetree.Node) interface{} {
				return n.Placements
			},
			want: []*move.Move{
				move.New(color.Black, point.New(0, 1)),
				move.New(color.Black, point.New(0, 2)),
				move.New(color.White, point.New(1, 1)),
				move.New(color.White, point.New(1, 2)),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := sgf.Parse(tc.sgf)
			if err != nil {
				t.Error(err)
				return
			}
			tp, err := movetree.ParsePath(tc.path)
			if err != nil {
				t.Error(err)
				return
			}
			got := tc.getter(tp.Apply(g.Root))
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("from node-getter and path %q, got %v, but wanted %v", tc.path, got, tc.want)
			}
		})
	}
}
