// Package kataprob finds problems from games that have katago-analysis data
// attached. It is assumed that all games have the relevant katago analysis data
// attached before this point.
package kataprob

import (
	"github.com/golang/glog"
	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/movetree"
	"github.com/otrego/clamshell/katago"
)

type FindBlunderOptions struct {
	PointThreshold float64
}

// FindBlunders finds positions (paths) that result from big swings in points.
func FindBlunders(g *movetree.MoveTree) ([]movetree.Path, error) {
	return findBlunders(g, &FindBlunderOptions{PointThreshold: 3.0})
}

func FindBlundersWithOptions(g *movetree.MoveTree, opt *FindBlunderOptions) ([]movetree.Path, error) {
	return findBlunders(g, opt)
}

func findBlunders(g *movetree.MoveTree, opt *FindBlunderOptions) ([]movetree.Path, error) {
	glog.V(3).Infof("Finding blunders with options: %v:\n", opt)
	var cur movetree.Path
	var found []movetree.Path
	if g.Root == nil {
		return found, nil
	}

	var prevLead float64

	for n := g.Root; n != nil; n = n.Next(0) {
		glog.V(3).Infof("VarNum %v\n", n.VarNum())
		glog.V(3).Infof("MoveNum %v\n", n.MoveNum())

		// We assume alternating moves. Lead is always presented as
		cur = append(cur, n.VarNum())
		glog.V(3).Infof("PrevLead %v\n", prevLead)

		d := n.AnalysisData()
		if d == nil {
			glog.Infof("nil analysis data")
			continue
		}

		katad, ok := d.(*katago.AnalysisResult)
		if !ok {
			glog.V(2).Infof("not analysisResult")
			continue
		}
		if katad.RootInfo == nil {
			// This
			glog.Errorf("no RootInfo for at move %v", n.MoveNum())
			continue
		}

		lead := katad.RootInfo.ScoreLead
		glog.V(3).Infof("Next ScoreLead: %v:", lead)

		// A positive ScoreLead means black is winning. Negative means white is winning.
		delta := lead - prevLead
		if n.Move.Color() == color.White {
			delta *= -1
		}
		glog.V(3).Infof("Delta: %v:", delta)

		if delta <= -opt.PointThreshold {
			found = append(found, cur.Clone())
		}

		prevLead = lead
	}

	return found, nil
}
