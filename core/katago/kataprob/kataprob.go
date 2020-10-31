// Package kataprob finds problems from games that have katago-analysis data
// attached. It is assumed that all games have the relevant katago analysis data
// attached before this point.
package kataprob

import (
	"github.com/otrego/clamshell/core/game"
	"github.com/otrego/clamshell/core/treepath"
)

// FindBlunders finds positions (paths) that result from big swings in points.
func FindBlunders(g *game.Game) ([]treepath.Treepath, error) {
	return nil, nil
}
