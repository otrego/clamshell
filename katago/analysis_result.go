package katago

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/otrego/clamshell/go/movetree"
)

// AnalysisResult represents the result of an analysis from katago.
//
// For more details, see: https://github.com/lightvector/KataGo/blob/master/docs/Analysis_Engine.md
type AnalysisResult struct {
	ID         string      `json:"id"`
	TurnNumber int         `json:"turnNumber"`
	MoveInfos  []*MoveInfo `json:"moveInfos"`
	RootInfo   *RootInfo   `json:"rootInfo"`

	// Not yet supported:
	// ownership
	// policy
}

// String returns the string form of the analysis result
func (ar *AnalysisResult) String() string {
	return fmt.Sprintf("{%s, %d, %v, %v}", ar.ID, ar.TurnNumber, ar.MoveInfos, ar.RootInfo)
}

// AnalysisList  is a collection of analysis results. This is normally what Katago
// produces when asked to do analysis.
type AnalysisList []*AnalysisResult

// ParseAnalysisList parses raw Katago analysis into an in-memory result list.
// We assume a collection of unordered AnalysisResult objects. Once parsed, the
// list is sorted by TurnNumber for convenience.
//
// Example:
//
//     {"id": "abcd", "turnNumber": 1, "moveInfos": [...]}
//     {"id": "abcd", "turnNumber"" 2, "moveInfos": [...]}
func ParseAnalysisList(content []byte) (AnalysisList, error) {
	dec := json.NewDecoder(strings.NewReader(string(content)))
	var out AnalysisList
	for {
		res := &AnalysisResult{}
		err := dec.Decode(res)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}
	out.Sort()
	return out, nil
}

// ParseAnalysis parses a single AnalysisResult from String
func ParseAnalysis(content string) (*AnalysisResult, error) {
	dec := json.NewDecoder(strings.NewReader(content))
	res := &AnalysisResult{}
	err := dec.Decode(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Sort the AnalysisList in-place based on TurnNumber.
func (al AnalysisList) Sort() {
	sort.SliceStable(al, func(i, j int) bool {
		return al[i].TurnNumber < al[j].TurnNumber
	})
}

// AddToGame attaches an analysis list to an existing movetree, based on turn
// number.
func (al AnalysisList) AddToGame(g *movetree.MoveTree) error {
	// make a shallow copy so we can sort without modifying recevier.
	alc := al[:]
	alc.Sort()

	if len(alc) == 0 {
		// Degenerate case, but not an error per-se.
		return errors.New("during AddToGame, no analysis data")
	}

	// Here we assume the AnalysisList is sorted.
	curNode := g.Root
	if curNode == nil {
		// Degenerate case.
		return errors.New("during AddToGame, nil root node")
	}

	anIdx := 0
	curAn := alc[anIdx]
	for {
		if curAn.TurnNumber == curNode.MoveNum() {
			// Match! Attach the analysis data.
			// Increment both the node and the analysis list.
			curNode.SetAnalysisData(curAn)
			if next := curNode.Next(0); next != nil {
				curNode = next
			} else {
				break
			}
			anIdx++
			if anIdx < len(alc) {
				curAn = alc[anIdx]
			} else {
				break
			}
		} else if curAn.TurnNumber > curNode.MoveNum() {
			// The analysis can get ahead if there are gaps (i.e., the move numbers
			// are not covering.
			if next := curNode.Next(0); next != nil {
				curNode = next
			} else {
				break
			}
		} else if curNode.MoveNum() > curAn.TurnNumber {
			// This shouldn't happen if we always start at the root of the movetree.
			// Return an error if it does.
			return fmt.Errorf("analysis TurnNumber %d got behind of the movetree move number %d", curAn.TurnNumber, curNode.MoveNum())
		}
	}
	return nil
}

// MoveInfo contains information about suggested moves.
type MoveInfo struct {
	// Move is a GTP point, and the move being analyzed
	Move string `json:"move"`

	// Visits is the number of visits invested into the move.
	Visits int `json:"visits"`

	// Winrate is the winrate of the move, as a float in [0,1].
	Winrate float64 `json:"winrate"`

	// StoreStdev is the predicted standard deviation of the final score of the
	// game after this move, in points. (NOTE: due to the mechanics of MCTS, this
	// value will be significantly biased high currently, although it can still be
	// informative as relative indicator).
	ScoreStdev float64 `json:"scoreStdev"`

	// ScoreLead is the predicted average number of points that the current side
	// is leading by (with this many points fewer, it would be an even game)
	ScoreLead float64 `json:"scoreLead"`

	// ScoreSelfPlay is the predicted average value of the final score of the game
	// after this move during selfplay, in points. (NOTE: users should usually
	// prefer scoreLead, since scoreSelfplay may be biased by the fact that KataGo
	// isn't perfectly score-maximizing).
	ScoreSelfPlay float64 `json:"scoreSelfPlay"`

	// Prior is the policy prior of the move, as a float in [0,1].
	Prior float64 `json:"prior"`

	// Utility of the move, combining both winrate and score, as a float in [-C,C]
	// where C is the maximum possible utility.
	Utility float64 `json:"utility"`

	// LCB is the move's winrate, as a float in [0,1].
	LCB float64 `json:"lcb"`

	// UtilityLCB is the lcb of the move's utility
	UtilityLCB float64 `json:"utilityLcb"`

	// Order is Katago's 's ranking of the move. 0 is the best, 1 is the next
	// best, and so on.
	Order int `json:"order"`

	// PV is the principal variation following this move. May be of variable
	// length or even empty.
	PV []string `json:"pv"`

	// Currently unsupported:
	// pvVisits
}

// RootInfo contains information for the move-at root.
type RootInfo struct {
	Winrate       float64 `json:"winrate"`
	ScoreLead     float64 `json:"scoreLead"`
	ScoreSelfPlay float64 `json:"scoreSelfPlay"`
	Utility       float64 `json:"utility"`
	Visits        int     `json:"visits"`
}
