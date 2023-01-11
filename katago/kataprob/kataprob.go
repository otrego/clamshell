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

// FindBlunderOptions contains options for blunder detection
type FindBlunderOptions struct {
	// Sets the minimum point delta beyond which a move is considered a blunder
	PointThreshold float64
	// Determines which color's blunders we are searching. Can be Empty as well
	Color color.Color
}

// FindBlunders finds positions (paths) that result from big swings in points.
func FindBlunders(g *movetree.MoveTree) ([]movetree.Path, error) {
	return findBlunders(g, &FindBlunderOptions{PointThreshold: 3.0, Color: color.Empty})
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
		if n.Move != nil {
			glog.V(3).Infof("Move %v\n", n.Move.GoString())
		}

		delta, ok := computeDelta(n, n.Parent)
		glog.V(3).Infof("Delta: %f\n", delta)
		if !ok {
			glog.V(3).Info("No ScoreLead for current node, skipping")
			continue
		}
		cur = append(cur, n.VarNum())

		if delta <= -opt.PointThreshold {
			found = append(found, cur.Clone())
			glog.V(3).Infof("Added to paths: %#v\n", found)
		}
	}

	if opt.Color != color.Empty {
		return filterColor(g, &found, opt.Color), nil
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

func filterColor(g *movetree.MoveTree, paths *[]movetree.Path, c color.Color) []movetree.Path {
	var filtered []movetree.Path
	for _, p := range *paths {
		if p.Apply(g.Root).Move.Color() == c {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
