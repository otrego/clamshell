package parser_test

import (
	"github.com/otrego/clamshell/core/parser"
	"reflect"
	"testing"
)

var sgfs = map[string]struct {
	raw string
	gm  string
	aw  []string
	ab  []string
	cm  []string
	b1  string
}{
	"descriptionTest": {
		raw: "(;GM[1]C[Try these Problems out!])",
		gm:  "1",
		cm:  []string{"Try these Problems out!"},
	},
	"base": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[0.00]
PW[White]PB[Black])`,
		gm: "1",
	},
	"escapedComment": {
		raw: "(;GM[1]FF[4]C[Josh[1k\\]: Go is Awesome!])",
		gm:  "1",
		cm:  []string{"Josh[1k\\]: Go is Awesome!"},
	},
	"veryEasy": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
C[Here's a basic example problem]RU[Japanese]SZ[19]KM[0.00]PW[j]PB[j]AW[ef]
;B[pd]
;W[cc]
;B[qf]
;W[nc]
;B[dd]
;W[pb])`,
		gm: "1",
		aw: []string{"ef"},
		cm: []string{"Here's a basic example problem"},
		b1: "pd",
	},
	"easy": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
 RU[Japanese]SZ[19]KM[0.00]
 PW[White]PB[Black]AW[pa][pb][sb][pc][qc][sc][qd][rd][sd]AB[oa][qa][ob][rb][oc][rc][pd][pe][qe][re][se]C[\\] Black to Live]
 (;B[sa];W[ra]C[Ko])
 (;B[ra]C[Correct];W[]C[And if white thinks it is seki?]
   (;B[qb]C[Correct.];W[sa];B[rb]C[Black lives])
   (;B[sa];W[qb];B[ra];W[rb]C[White Lives])
 )
 (;B[qb];W[ra]C[White lives]))`,
		gm: "1",
		aw: []string{"pa", "pb", "sb", "pc", "qc", "sc", "qd", "rd", "sd"},
		ab: []string{"oa", "qa", "ob", "rb", "oc", "rc", "pd", "pe", "qe", "re", "se"},
		cm: []string{"\\\\] Black to Live"},
		b1: "sa",
	},
	"marky": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[0.00]
PW[White]PB[Black]CR[rb][rc][re]AB[pc][qd][pe]
[re]LB[pb:3][qb:2][pc:B][qc:1][pd:A]TR[qd][qe]SQ[rd:re]
;B[sa]TR[qa]C[bar]
;W[fi]SQ[ab]C[foo])`,
		gm: "1",
		ab: []string{"pc", "qd", "pe", "re"},
		b1: "sa",
	},
	"trivialProblem": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[0.00]
PW[White]PB[Black]GB[1]
C[Here's an example diagram. I have marked 1 on the diagram.
Let's pretend it was white's last move.  Think on this move, since
it may be a problem in the near future!]
LB[pb:1]
AW[pb][mc][pc][qd][rd][qf][pg][qg]
AB[jc][oc][qc][pd][pe][pf])`,
		gm: "1",
		aw: []string{"pb", "mc", "pc", "qd", "rd", "qf", "pg", "qg"},
		ab: []string{"jc", "oc", "qc", "pd", "pe", "pf"},
		cm: []string{"Here's an example diagram. I have marked 1 on the diagram.\nLet's pretend it was white's last move.  Think on this move, since\nit may be a problem in the near future!"},
	},
	"realProblem": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[0.00]
PW[White]PB[Black]AW[pb][mc][pc][qd][rd][qf][pg][qg]
AB[jc][oc][qc][pd][pe][pf]
C[Look Familiar?]
(;B[ob]C[Bad style.]
;W[qb]
(;B[nd]C[White's stone can easily escape.])
(;B[me]C[Lots of bad aji.]))
(;B[nc]
(;W[qb]
;B[md]C[Correct]GB[1])
(;W[md]
;B[qb]GB[1]C[White loses his corner])))`,
		gm: "1",
		aw: []string{"pb", "mc", "pc", "qd", "rd", "qf", "pg", "qg"},
		ab: []string{"jc", "oc", "qc", "pd", "pe", "pf"},
		cm: []string{"Look Familiar?"},
		b1: "ob",
	},
	"complexProblem": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[Glift]ST[2]
RU[Japanese]SZ[19]KM[0.00]
C[Black to play. There aren't many option
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
		gm: "1",
		aw: []string{"pa", "qa", "nb", "ob", "qb", "oc", "pc", "md", "pd", "ne", "oe"},
		ab: []string{"na", "ra", "mb", "rb", "lc", "qc", "ld", "od", "qd", "le", "pe", "qe", "mf", "nf", "of", "pg"},
		cm: []string{"Black to play. There aren't many option\nto choose from, but you might be surprised at the answer!"},
		b1: "mc",
	},
	"markTest": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[0.00]
C[[Mark Test\]]
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
		gm: "1",
		aw: []string{"na", "oa", "pa", "qa", "ra", "sa", "ka", "la", "ma", "ja"},
		ab: []string{"nb", "ob", "pb", "qb", "rb", "sb", "kb", "lb", "mb", "jb"},
		cm: []string{"[Mark Test\\]"},
	},
	"twoOptions": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[0.00]
PW[White]PB[Black]EV[ALL_CORRECT]AW[oc][pe]AB[mc][qd]C[What are the normal ways black follows up this position?]
(;B[pd]C[Correct]
;W[od]
;B[oe])
(;B[qe]C[Correct]
;W[pf]
;B[qg]))`,
		gm: "1",
		aw: []string{"oc", "pe"},
		ab: []string{"mc", "qd"},
		cm: []string{"What are the normal ways black follows up this position?"},
		b1: "pd",
	},
	"passingExample": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[0.00]
PW[White]PB[Black]
;B[]
;AW[qc]AB[pd]C[How should White respond?]
(;W[pc]
;B[od]C[Correct])
(;W[qd]
;B[pe]C[Correct]))`,
		gm: "1",
	},
	"goGameGuruHard": {
		raw: `
(;GM[2]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[0.00]C[A Problem from GoGameGuru]
AW[po][qo][ro][so][np][op][pq][nr][pr][qr][rs]
AB[qm][on][pn][oo][pp][qp][rp][sp][qq][rr][qs]
(;B[sr]
(;W[rq];B[sq];W[ps];B[rn]C[Correct])
(;W[ps]
(;B[rn];W[rq];B[sq];W[qs]
(;B[sn]C[Correct])
(;B[qn]C[Correct]))
(;B[qn];W[rq];B[sq];W[qs];B[rn]C[Correct])
(;B[sn];W[rq];B[sq];W[qs];B[rn]C[Correct])))
(;B[sq];W[ps]
(;B[rn];W[sr];B[ss]C[It's a ko, but black can do better.])
(;B[sr];W[qs];B[rn];W[ss])
(;B[qn];W[sr];B[ss]C[It's a ko, but black can do better.])
(;B[sn];W[sr];B[ss]C[It's a ko, but black can do better.]))
(;B[ss];W[sq];B[rq];W[ps]
(;B[rn];W[rs]C[It's a ko, but black can do better.])
(;B[qn];W[rs]C[It's a ko, but black can do better.])
(;B[sn];W[rs]C[It's a ko, but black can do better.]))
(;B[rq];W[ps]
(;B[sr];W[qs]
(;B[rn];W[ss])
(;B[qn];W[ss]))
(;B[rn];W[sr])
(;B[qn];W[sr]))
(;B[rn];W[sq])
(;B[qn];W[sq])
(;B[sn];W[sq]))`,
		gm: "2",
		aw: []string{"po", "qo", "ro", "so", "np", "op", "pq", "nr", "pr", "qr", "rs"},
		ab: []string{"qm", "on", "pn", "oo", "pp", "qp", "rp", "sp", "qq", "rr", "qs"},
		cm: []string{"A Problem from GoGameGuru"},
		b1: "sr",
	},
	"leeGuGame6": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Chinese]SZ[19]KM[7.50]TM[14100]OT[5x60 byo-yomi]
GN[Lee-Sedol-vs-Gu-Li-20140525]PW[Lee Sedol]PB[Gu
Li]WR[9d]BR[9d]DT[2014-07-27]EV[MLily Gu vs Lee Jubango]RO[Game
6]PC[Liuan, Anhui, China]SO[http://gogameguru.com/]RE[W+Resign] ;B[pd]
;W[dp] ;B[qp] ;W[dc] ;B[oq] ;W[qf] ;B[pi] ;W[nd] ;B[pf] ;W[qc] ;B[pc]
;W[pg] ;B[pe] ;W[ph] ;B[qi] ;W[ng] ;B[re] ;W[ni] ;B[rg] ;W[ce] ;B[jd]
;W[hc] ;B[ne] ;W[nk] ;B[mh] ;W[nh] ;B[me] ;W[fq] ;B[kp] ;W[pk] ;B[jb]
;W[dj] ;B[cn] ;W[bp] ;B[bj] ;W[jp] ;B[jq] ;W[cj] ;B[bk] ;W[bi] ;B[dl]
;W[iq] ;B[ir] ;W[ip] ;B[hr] ;W[kq] ;B[jr] ;W[lq] ;B[lp] ;W[mp] ;B[mo]
;W[np] ;B[mq] ;W[nq] ;B[mr] ;W[no] ;B[ln] ;W[mn] ;B[lo] ;W[lm] ;B[lr]
;W[el] ;B[em] ;W[cl] ;B[ep] ;W[eq] ;B[ek] ;W[fl] ;B[dk] ;W[ck] ;B[fk]
;W[fi] ;B[do] ;W[cp] ;B[mm] ;W[nn] ;B[ml] ;W[om] ;B[km] ;W[ll] ;B[mk]
;W[kl] ;B[lj] ;W[qn] ;B[ol] ;W[ok] ;B[pl] ;W[ql] ;B[nl] ;W[rk] ;B[qj]
;W[qk] ;B[pm] ;W[pn] ;B[ro] ;W[rn] ;B[jk] ;W[il] ;B[ik] ;W[hk] ;B[jl]
;W[gk] ;B[fm] ;W[kn] ;B[kr] ;W[jm] ;B[hm] ;W[gl] ;B[gm] ;W[im] ;B[hl]
;W[hj] ;B[gi] ;W[fj] ;B[in] ;W[kk] ;B[ii] ;W[kj] ;B[ij] ;W[ki] ;B[gj]
;W[lh] ;B[fc] ;W[ic] ;B[dd] ;W[cd] ;B[db] ;W[cc] ;B[ec] ;W[jc] ;B[kc]
;W[kd] ;B[lc] ;W[fe] ;B[cb] ;W[bb] ;B[gd] ;W[ge] ;B[hd] ;W[id] ;B[he]
;W[je] ;B[hf] ;W[gg] ;B[gh] ;W[hg] ;B[ig] ;W[if] ;B[jg] ;W[kf] ;B[fg]
;W[gf] ;B[de] ;W[ef] ;B[df] ;W[eg] ;B[ba] ;W[ab] ;B[gb] ;W[hb] ;B[ga]
;W[cf] ;B[dg] ;W[eh] ;B[kg] ;W[lf] ;B[dh] ;W[ed] ;B[lg] ;W[mg] ;B[mf]
;W[ee] ;B[le] ;W[gc] ;B[fd] ;W[ih] ;B[kh] ;W[jh]C[http://gogameguru.com/])`,
		gm: "1",
		b1: "pd",
	},
	"yearbookExample": {
		raw: `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[6.50]
PW[Lee Sedol]PB[Gu Li]WR[9d]BR[7d]DT[2004-11-16]EV[9th Samsung Cup]RO[Semifinal]PC[Ulsan]SO[https://gogameguru.com/]RE[W+Resign]
;B[qd] ;W[dd] ;B[pq] ;W[oc] ;B[dp] ;W[po] ;B[pe] ;W[md] ;B[qm] ;W[qq] ;B[qp]
;W[pp] ;B[qo] ;W[qn] ;B[pn] ;W[rn] ;B[rq] ;W[qr] ;B[ro] ;W[rm] ;B[oq] ;W[np]
;B[rr] ;W[ql] ;B[pm] ;W[pl] ;B[nm] ;W[op] ;B[ol] ;W[pj] ;B[qh] ;W[ok] ;B[nk]
;W[nj] ;B[mk] ;W[so] ;B[rp] ;W[mm] ;B[nn] ;W[mn] ;B[mj] ;W[ni] ;B[mi] ;W[mh]
;B[lh] ;W[mg] ;B[lg] ;W[on] ;B[om] ;W[mf] ;B[jp] ;W[km] ;B[jj] ;W[im] ;B[lp]
;W[nq] ;B[pr] ;W[or] ;B[qs] ;W[no] ;B[nl] ;W[lo] ;B[gp] ;W[jh] ;B[ji] ;W[ih]
;B[hj] ;W[fm] ;B[kf] ;W[if] ;B[kd] ;W[id] ;B[fj] ;W[dm] ;B[ck] ;W[kk] ;B[le]
;W[of] ;B[cf] ;W[dh] ;B[hg] ;W[hh] ;B[gh] ;W[gg] ;B[fh] ;W[hi] ;B[ik] ;W[gi]
;B[fi] ;W[gj] ;B[gk] ;W[fk] ;B[gl] ;W[fl] ;B[gm] ;W[gn] ;B[hn] ;W[hm] ;B[fn]
;W[go] ;B[fo] ;W[ho] ;B[il] ;W[ej] ;B[jn] ;W[jm] ;B[in] ;W[hp] ;B[gq] ;W[jo]
;B[io] ;W[ip] ;B[ko] ;W[kn] ;B[jo] ;W[jq] ;B[kp] ;W[hr] ;B[kr] ;W[gr] ;B[jr]
;W[fp] ;B[fq] ;W[fr] ;B[eq] ;W[ir] ;B[er] ;W[lr] ;B[mr] ;W[lq] ;B[kq] ;W[mq]
;B[iq] ;W[hq] ;B[fs] ;W[ks] ;B[is] ;W[mp] ;B[fg] ;W[gf] ;B[ff] ;W[en] ;B[eo]
;W[do] ;B[ep] ;W[co] ;B[gs] ;W[jq] ;B[ge] ;W[hf] ;B[iq] ;W[hl] ;B[hk] ;W[jq]
;B[jg] ;W[js] ;B[ig] ;W[cq] ;B[cp] ;W[bp] ;B[br] ;W[fd] ;B[cd] ;W[cc] ;B[de]
;W[bd] ;B[dl] ;W[ek] ;B[cm] ;W[bn] ;B[bq] ;W[cj] ;B[bj] ;W[ci] ;B[bm] ;W[bi]
;B[ed] ;W[ce] ;B[dc] ;W[ee] ;B[cd] ;W[jf] ;B[kg] ;W[dd] ;B[ei] ;W[dj] ;B[cd]
;W[cr] ;B[ap] ;W[dd] ;B[df] ;W[ec] ;B[dg] ;W[bf] ;B[bg] ;W[ag] ;B[ah] ;W[af]
;B[fe] ;W[bh] ;B[gd] ;W[qg] ;B[rg] ;W[qf] ;B[rf] ;W[qc] ;B[rk] ;W[rl] ;B[qi]
;W[rj] ;B[qj] ;W[qk] ;B[oh] ;W[og] ;B[oi] ;W[oj] ;B[ri] ;W[sk] ;B[rd] ;W[rc]
;B[od] ;W[pc] ;B[me] ;W[nh] ;B[si] ;W[ao] ;B[bo]
;W[sd]C[### 9th Samsung Cup - Game One
**November 16 2004, Ulsan, Korea: 9th Samsung Cup Semifinals, Game One**
*Gu Li 7d (black) vs Lee Sedol 9d*
White 38 was painful for Black.
White 76 was a mistake.
Black 125 was the losing move. Black should have played Black 125 at Black 145 instead.
**228 moves: White won by resignation.**])`,
		gm: "1",
		b1: "qd",
	},
}

func TestFoo(t *testing.T) {
	for name, tt := range sgfs {
		t.Run(tt.raw, func(t *testing.T) {
			p := parser.FromString(tt.raw)
			game, err := p.Parse()
			if err != nil {
				t.Error(err)
			}
			node := game.Root.Children[0]

			if gm := node.Properties["GM"][0]; gm != tt.gm {
				t.Errorf("%s: GM=%v, expected %v", name, gm, tt.gm)
			}

			if aw := node.Properties["AW"]; !reflect.DeepEqual(aw, tt.aw) {
				t.Errorf("%s: AW=%v, expected %v", name, aw, tt.aw)
			}

			if ab := node.Properties["AB"]; !reflect.DeepEqual(ab, tt.ab) {
				t.Errorf("%s: AB=%v, expected %v", name, ab, tt.ab)
			}

			if cm := node.Properties["C"]; !reflect.DeepEqual(cm, tt.cm) {
				t.Errorf("%s: C=%v, expected %v", name, cm, tt.cm)
			}

			if len(node.Children) > 0 {
				nxt := node.Children[0]
				if b1 := nxt.Properties["B"][0]; b1 != tt.b1 {
					t.Errorf("%s: B=%v, expected %v", name, b1, tt.b1)
				}
			}
		})
	}
}
