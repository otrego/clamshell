package sgf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/sgf"
)

func TestSerialize(t *testing.T) {
	testCases := []struct {
		desc string
		sgf  string
	}{
		{
			desc: "black move",
			sgf:  "(;GM[1];B[ab])",
		},
		{
			desc: "white move",
			sgf:  "(;GM[1];W[ab])",
		},
		{
			desc: "black & white placements",
			sgf:  "(;GM[1];AB[ab][ac]AW[bb][bc])",
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
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			//Create game
			g, err := sgf.Parse(tc.sgf)
			if err != nil {
				t.Error(err)
				return
			}

			//convert original movetree back to sgf, then back to new game
			serialized, err := sgf.Serialize(g)
			if err != nil {
				t.Fatal(err)
			}
			got, err := sgf.Parse(serialized)

			//check if both movetrees are identical
			var sbWant strings.Builder
			var sbGot strings.Builder
			g.Root.Traverse(func(n *movetree.Node) {
				sbWant.WriteString(fmt.Sprintf("%v", n.SGFProperties))
			})

			got.Root.Traverse(func(n *movetree.Node) {
				sbGot.WriteString(fmt.Sprintf("%v", n.SGFProperties))
			})

			if sbWant.String() != sbGot.String() {
				t.Errorf("want:\n%s\ngot%s", sbWant.String(), sbGot.String())
			}
		})
	}
}
