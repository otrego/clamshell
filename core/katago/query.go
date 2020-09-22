package katago

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/otrego/clamshell/core/game"
	"github.com/otrego/clamshell/core/point"
)

// Query is a katago analysis query. Here we use pointer values for primitives
// when types are optional.
type Query struct {
	// ID is a *required* arbitrary identifier for the analysis.
	ID string `json:"id"`

	// InitialStones are the initial stones on the board.
	InitialStones []Move `json:"initialStones,omitempty"`

	// InitialPlayer is the initial player
	InitialPlayer string `json:"initialPlayer,omitempty"`

	// Moves are the moves played in the game.
	Moves []Move `json:"moves"`

	// Rules are a *required* string shorthand for the rules. Technically this can
	// either be a string or a full
	Rules Rules `json:"rules"`

	// Komi is the komi to be used. Generally this is implied by the rules, so
	// doesn't need to be set. Note that the komi should be integers or half
	// integers: ex: 8, 7.5, 6.5 etc.
	Komi *float64 `json:"komi,omitempty"`

	// BoardXSize is the width of the board
	BoardXSize *int `json:"boardXSize,omitempty"`

	// BoardXSize is the width of the board
	BoardYSize *int `json:"boardYSize,omitempty"`

	// AnalyzeTurns is the turns of the game to analyze. If this field is not
	// specified, defaults to analyzing the last turn only.
	AnalyzeTurns []int `json:"analyzeTurns,omitempty"`

	// MaxVisits is the maximum number of visits to use. If not specified,
	// defaults to the value in the analysis config file. If specified, overrides
	// it.
	MaxVisits *int `json:"maxVisits,omitempty"`

	// Not yet supported options
	// See: https://github.com/lightvector/KataGo/blob/master/docs/Analysis_Engine.md
	// whiteHandicapBonus
	// rootPolicyTemperature
	// rootFpuReductionMax
	// includeOwnership
	// includePolicy
	// includePVVisits
	// avoidMoves
	// allowMoves
	// overrideSettings
	// priority
}

// NewQuery creates an analysis Query object with default parameters
func NewQuery() *Query {
	return &Query{
		ID:    uuid.New().String(),
		Rules: TrompTaylorRules,
	}
}

// Move is a string-tuple having the form [<MOVE>, <POS>].
// For example : ["B", "C4"]
type Move [2]string

// Rules are rule-alias strings.
type Rules string

const (
	// TrompTaylorRules are modern area scoring rules.
	TrompTaylorRules Rules = "tromp-taylor"
	// ChineseRules are traditional chinese rules with simple ko.
	ChineseRules Rules = "chinese"
	// ChineseOGSRules are chinese rules as implemented in OGS/KGS.
	ChineseOGSRules Rules = "chinese-ogs"
	// NewZealandRules are new-zealand rules, which also use area scoring.
	NewZealandRules Rules = "new-zealand"
	// Japanese Rules (equiv to korean) are traditional territory scoring rules.
	Japanese Rules = "japanese"
)

// ToJSON converts a query to JSON.
func (q *Query) ToJSON() ([]byte, error) {
	return json.Marshal(q)
}

type gameConverter struct {
	g *game.Game
}

func (gc *gameConverter) point(pt *point.Point) string {
	const a = 'A'
	val := rune(a + pt.X())
	return fmt.Sprintf("%c%d", val, pt.Y()+1)
}

func (gc *gameConverter) move(mv *game.Move) Move {
	return Move{string(mv.Color()), gc.point(mv.Point())}
}

func (gc *gameConverter) initialStones() []Move {
	var out []Move
	for _, mv := range gc.g.Root.Placements {
		out = append(out, gc.move(mv))
	}
	return out
}

func (gc *gameConverter) initialPlayer() string {
	if val, ok := gc.g.Root.Properties["PL"]; ok && len(val) > 0 {
		return val[0]
	}
	return ""
}

func (gc *gameConverter) mainBranchMoves(maxMoves int) []Move {
	var out []Move
	idx := 0
	for n := gc.g.Root; ; n = n.Children[0] {
		if n.Move != nil && !n.Move.IsPass() {
			out = append(out, gc.move(n.Move))
		}
		if len(n.Children) == 0 {
			// No more children; terminate traversal.
			break
		}
		if idx >= maxMoves {
			break
		}
		idx++
	}
	return out
}

func (gc *gameConverter) rules() Rules {
	if val, ok := gc.g.Root.Properties["RU"]; ok && len(val) > 0 {
		return Rules(val[0])
	}
	return TrompTaylorRules
}

func (gc *gameConverter) komi() (*float64, error) {
	if val, ok := gc.g.Root.Properties["KM"]; ok && len(val) > 0 {
		km, err := strconv.ParseFloat(val[0], 64)
		if err != nil {
			return nil, err
		}
		return &km, nil
	}
	return nil, nil
}

func (gc *gameConverter) boardSize() (*int, error) {
	if val, ok := gc.g.Root.Properties["SZ"]; ok && len(val) > 0 {
		sz, err := strconv.Atoi(val[0])
		if err != nil {
			return nil, err
		}
		if sz > 25 {
			return nil, fmt.Errorf("only sizes up to 25 are supported")
		}
		return &sz, nil
	}
	return nil, nil
}

func (gc *gameConverter) analyzeMainBranch(maxMoves int) []int {
	var out []int
	idx := 0
	for n := gc.g.Root; ; n = n.Children[0] {
		if n.Move != nil && !n.Move.IsPass() {
			out = append(out, idx)
		} else if n.Move != nil && n.Move.IsPass() {
			// stop analyzing at first pass
			break
		}
		if len(n.Children) == 0 {
			// No more children; terminate traversal.
			break
		}
		if idx >= maxMoves {
			break
		}
		idx++
	}
	return out
}

// MainBranchSurvey does a full game analysis of the main-branch.
func MainBranchSurvey(g *game.Game) (*Query, error) {
	q := NewQuery()
	gc := &gameConverter{g: g}

	maxMoves := 10

	q.InitialStones = gc.initialStones()
	q.InitialPlayer = gc.initialPlayer()
	q.Moves = gc.mainBranchMoves(maxMoves)
	q.Rules = gc.rules()
	km, err := gc.komi()
	if err != nil {
		return nil, err
	}
	q.Komi = km

	sz, err := gc.boardSize()
	if err != nil {
		return nil, err
	}
	q.BoardYSize = sz
	q.BoardXSize = sz
	// q.AnalyzeTurns = gc.analyzeMainBranch(maxMoves)
	// q.AnalyzeTurns = gc.analyzeMainBranch(maxMoves)

	return q, nil
}
