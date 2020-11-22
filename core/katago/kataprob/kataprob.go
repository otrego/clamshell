// Package kataprob finds problems from games that have katago-analysis data
// attached. It is assumed that all games have the relevant katago analysis data
// attached before this point.
package kataprob

import (
	"math"

	"github.com/golang/glog"
	"github.com/otrego/clamshell/core/game"
	"github.com/otrego/clamshell/core/katago"
	"github.com/otrego/clamshell/core/treepath"
)

// FindBlunders finds positions (paths) that result from big swings in points.
func FindBlunders(g *game.Game) ([]treepath.Treepath, error) {
	blunderAmt := 3.0

	var cur treepath.Treepath
	var found []treepath.Treepath
	if g.Root == nil {
		return found, nil
	}

	var prevLead float64

	for n := g.Root; n != nil; n = n.Next(0) {
		glog.Infof("VarNum %v\n", n.VarNum())
		glog.Infof("MoveNum %v\n", n.MoveNum())

		// We assume alternating moves. Lead is always presented as
		pl := prevLead
		cur = append(cur, n.VarNum())
		glog.Infof("PrevLead %v\n", prevLead)

		d := n.AnalysisData()
		if d == nil {
			glog.Infof("nil analysis data")
			continue
		}

		katad, ok := d.(*katago.AnalysisResult)
		if !ok {
			glog.Infof("not analysisResult")
			continue
		}
		if katad.RootInfo == nil {
			glog.Infof("no RootInfo")
			continue
		}

		lead := katad.RootInfo.ScoreLead
		nextLead := -1 * lead
		glog.Infof("Next ScoreLead: %v:", nextLead)
		delta := nextLead - pl
		glog.Infof("Delta: %v:", delta)
		if delta >= math.Abs(blunderAmt) {
			found = append(found, cur.Clone())
		}

		// prevLead is always with respect to current player
		prevLead = nextLead
	}

	return found, nil
}
