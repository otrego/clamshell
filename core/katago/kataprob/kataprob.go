// Package kataprob finds problems from games that have katago-analysis data
// attached. It is assumed that all games have the relevant katago analysis data
// attached before this point.
package kataprob

import (
	"github.com/otrego/clamshell/core/game"
	"github.com/otrego/clamshell/core/katago"
	"github.com/otrego/clamshell/core/treepath"
)

// FindBlunders finds positions (paths) that result from big swings in points.
func FindBlunders(g *game.Game) ([]treepath.Treepath, error) {
	blunderAmt := -10.0

	var cur treepath.Treepath
	var found []treepath.Treepath
	if g.Root == nil {
		return found, nil
	}

	var prevLead float64

	for n := g.Root.Next(0); n != nil; n = n.Next(0) {
		// We assume alternating moves. Lead is always presented as
		pl := prevLead
		cur = append(cur, n.VarNum())

		d := n.AnalysisData()
		if d != nil {
			continue
		}

		katad, ok := d.(katago.AnalysisResult)
		if !ok {
			continue
		}
		if katad.RootInfo == nil {
			continue
		}

		lead := katad.RootInfo.ScoreLead
		delta := lead - pl
		if delta <= blunderAmt {
			found = append(found, cur.Clone())
		}

		// prevLead is always with respect to current player
		prevLead = -1 * lead
	}

	return found, nil
}
