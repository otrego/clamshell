package sgf

import (
	"fmt"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	s := `
(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[6.50]
PW[Player White [1k\]]PB[Player Black [3k\]]C[This (yes;) is a [5k\] comment]
;B[pd]C[Watcher [10k\]: hello world]
;W[dd]
;B[pq]
;W[dp]
(;B[qk]
(;W[fq]
;B[nc]
;W[fc]
;B[cj]
;W[cl]
;B[cf]
;W[mp]
;B[po]
;W[jp]
;B[ej])
(;W[nc]
;B[pf]
;W[pb]
(;B[qc]
;W[kc])
(;B[fq]))
(;W[mp]
;B[po]
;W[jp]))
(;B[fq]
(;W[cn]
;B[jp])
(;W[qo])))`
	r := strings.NewReader(s)
	parser := NewParser(r)
	game, err := parser.Parse()
	if err != nil {
		t.Error(err)
	}
	if pw := game.Nodes[0].Fields["PW"]; pw != "Player White [1k\\]" {
		t.Errorf("PW=%v, expected %v", pw, "Player White [1k\\]")
	}
	if pb := game.Nodes[0].Fields["PB"]; pb != "Player Black [3k\\]" {
		t.Errorf("PB=%v, expected %v", pb, "Player Black [3k\\]")
	}
	for i := 0; i < game.Index; i++ {
		fmt.Println(game.Nodes[i])
	}
}
