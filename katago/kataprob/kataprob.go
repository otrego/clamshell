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

	for n := g.Root; n != nil; n = n.Next(0) {
		glog.V(3).Infof("VarNum %v\n", n.VarNum())
		glog.V(3).Infof("MoveNum %v\n", n.MoveNum())

		cur = append(cur, n.VarNum())

		delta, ok := computeDelta(n, n.Parent)
		if !ok {
			glog.V(3).Info("No ScoreLead for current node, skipping")
			continue
		}

		if delta <= -opt.PointThreshold {
			found = append(found, cur.Clone())
		}
	}

	return found, nil
}

func getScoreLead(n *movetree.Node) (float64, bool) {
	if n == nil {
		glog.Info("nil node")
		return 0, false
	}

	d := n.AnalysisData()
	if d == nil {
		glog.Info("nil analysis data")
		return 0, false
	}

	nData, ok := d.(*katago.AnalysisResult)
	if !ok {
		glog.V(2).Info("not analysisResult")
		return 0, false
	}

	if nData.RootInfo == nil {
		glog.Errorf("no RootInfo for at move %v", n.MoveNum())
		return 0, false
	}

	return nData.RootInfo.ScoreLead, true
}

func computeDelta(n, p *movetree.Node) (float64, bool) {
	previousLead, ok := getScoreLead(p)
	if !ok {
		glog.V(3).Info("Defaulting score to 0")
	}
	glog.V(3).Infof("Previous Lead %v\n", previousLead)

	currentLead, ok := getScoreLead(n)
	if !ok {
		return 0, false
	}
	glog.V(3).Infof("Current Lead %v\n", currentLead)

	// A positive ScoreLead means black is winning. Negative means white is winning.
	delta := currentLead - previousLead
	if n.Move.Color() == color.White {
		delta *= -1
	}
	glog.V(3).Infof("Delta: %v:", delta)
	return delta, true
}
