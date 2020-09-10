package katago

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Query is a katago analysis query. Here we use pointer values for primitives
// when types are optional.
type Query struct {
	// ID is a *required* arbitrary identifier for the analysis.
	ID string `json:"id"`

	// InitialStones are the initial stones on the board.
	InitialStones []Move `json:"initialStones"`

	// Moves are the moves played in the game.
	Moves []Move `json:"moves"`

	// Rules are a *required* string shorthand for the rules. Technically this can
	// either be a string or a full
	Rules Rules `json:"rules"`

	// Komi is the komi to be used. Generally this is implied by the rules, so
	// doesn't need to be set.
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
type Move []string

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
