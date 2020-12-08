package problems_test

import (
	"testing"

	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/problems"
	"github.com/otrego/clamshell/core/prop"
	"github.com/otrego/clamshell/core/sgf"
)

func TestFlatten(t *testing.T) {

	testCases := []struct {
		desc string
		tp   string
		sgf  string
	}{
		{
			desc: "entire game flatten",
			tp:   "0:140",
			sgf: `(;FF[4]CA[UTF-8]GM[1]DT[2020-08-05]PB[player1]PW[player2]
				BR[2k] WR[5d] TM[259200]OT[86400 fischer] RE[W+R] SZ[19]
				KM[6.5] RU[Japanese] ;B[pd] (;W[dc] (;B[qp] (;W[cq] (;B[np]
				(;W[qf] (;B[nc] (;W[qd] (;B[qc] (;W[rc] (;B[qe] (;W[rd]
				(;B[re] (;W[pe] (;B[rf] (;W[pc] (;B[od] (;W[pb] (;B[qg]
				(;W[df] (;B[oh] (;W[lq] (;B[lo] (;W[jp] (;B[mq] (;W[lp]
				(;B[mo] (;W[qo] (;B[pp] (;W[rp] (;B[rq] (;W[pl]
				(;B[qm] (;W[ql] (;B[qn] (;W[qi] (;B[nl] (;W[nk]
				(;B[ok] (;W[ol] (;B[oj] (;W[nm] (;B[ml] (;W[mm] (;B[ll]
				(;W[lm] (;B[jo] (;W[io] (;B[jn] (;W[kl] (;B[mj] (;W[kk] (;B[lj]
				(;W[jm] (;B[ip] (;W[hp] (;B[iq] (;W[in] (;B[hq] (;W[gp]
				(;B[gq] (;W[fq] (;B[kp] (;W[dk] (;B[fp] (;W[fr] (;B[fo]
				(;W[gn] (;B[im] (;W[hn] (;B[jq] (;W[il] (;B[cp] (;W[dq]
				(;B[bk] (;W[ck] (;B[bl] (;W[cm] (;B[ci] (;W[bm] (;B[bf] 
				(;W[cf] (;B[bg] (;W[be] (;B[ai] (;W[ic] (;B[om] (;W[lc]
				(;B[hd] (;W[hc] (;B[db] (;W[cc] (;B[fd] (;W[ch] (;B[bh]
				(;W[ge] (;B[gd] (;W[je] (;B[gf] (;W[gg] (;B[hf] (;W[hg]
				(;B[if] (;W[kg] (;B[fg] (;W[fh] (;B[jh] (;W[ni] (;B[oi]
				(;W[kh] (;B[eh] (;W[eg] (;B[ff] (;W[gi] (;B[fb] (;W[ec] 
				(;B[eb] (;W[gb] (;B[fc] (;W[cb] (;B[jf] (;W[kf] (;B[ke]
				(;W[le] (;B[kd] (;W[ld] (;B[ji] (;W[ki] (;B[kj] (;W[jj]
				(;B[ij] (;W[jk] (;B[nj] (;W[ee] (;B[jd] (;W[id] (;B[ie]
				(;W[fe] )))))))))))))))))))))))))))))))))))))))))))))))
				)))))))))))))))))))))))))))))))))))))))))))))))))))))))
				))))))))))))))))))))))))))))))))))))`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tp, _ := movetree.ParsePath(tc.tp)
			g, _ := sgf.Parse(tc.sgf)
			bWant, _ := problems.PopulateBoard(tp, g)

			// Test that the boards are identical
			tpRoot, _ := movetree.ParsePath("0:2")
			gFlat, err := problems.Flatten(tp, g)
			if err != nil {
				t.Errorf("%e", err)
			}
			bGot, _ := problems.PopulateBoard(tpRoot, gFlat)
			if bWant.String() != bGot.String() {
				t.Errorf("wanted \n %s \n\n but got \n %s", bWant.String(), bGot.String())
			}

			// Test that the properties are identical
			propsWant, _ := prop.ConvertNode(g.Root)
			propsGot, _ := prop.ConvertNode(gFlat.Root)

			if bWant.String() != bGot.String() {
				t.Errorf("wanted %s \n but got %s", propsWant, propsGot)
			}
		})
	}
}
