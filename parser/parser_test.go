package parser_test

import (
	"github.com/otrego/clamshell/parser"
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
	p := parser.FromString(s)
	game, err := p.Parse()
	if err != nil {
		t.Error(err)
	}
	node := game.Root.Children[0]

	if pw := node.Properties["PW"][0]; pw != "Player White [1k\\]" {
		t.Errorf("PW=%v, expected %v", pw, "Player White [1k\\]")
	}
	if pb := node.Properties["PB"][0]; pb != "Player Black [3k\\]" {
		t.Errorf("PB=%v, expected %v", pb, "Player Black [3k\\]")
	}
}
