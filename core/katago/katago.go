// Package katago provides wrappers for analyzing games with katago.
package katago

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Analyzer is a katago-analyzer
type Analyzer struct {
	// Model is the path to the model file.
	Model string

	// Config is the path to the config file.
	Config string
}

// Query is a katago analysis query. Here we use pointer values for primitives
// when types are optional.
type Query struct {
	// Id is a *required* arbitrary identifier for the analysis.
	Id string `json:"id"`

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

	// Not supported options
	// See: https://github.com/lightvector/KataGo/blob/master/docs/Analysis_Engine.md
	//
	// whiteHandicapBonus (0|N|N-1): Optional. See kata-get-rules in GTP
	// Extensions for what these mean. Can be used to override the handling of
	// handicap bonus, taking precedence over rules.
	//
	// rootPolicyTemperature (float): Optional. Set this to a value > 1 to make
	// KataGo do a wider search.
	//
	// rootFpuReductionMax (float): Optional. Set this to 0 to make KataGo more
	// willing to try a variety of moves.
	//
	// includeOwnership (boolean): Optional. If true, report ownership prediction
	// as a result. Will double memory usage and reduce performance slightly.
	//
	// includePolicy (boolean): Optional. If true, report neural network raw
	// policy as a result. Will not signficiantly affect performance.
	//
	// includePVVisits (boolean): Optional. If true, report the number of visits
	// for each move in any reported pv.
	//
	// avoidMoves (list of dicts): Optional. If true, UNTILDEPTH` - Prohibit the
	// search from exploring the specified moves for the specified player, until a
	// certain number of ply deep in the search. Each dict must contain these
	// fields:
	//
	// allowMoves (list of dicts): Optional. Same as avoidMoves except prohibits
	// all moves EXCEPT the moves specified. Currently, the list of dicts must
	// also be length 1.
	//
	// overrideSettings (object): Optional. Specify any number of paramName:value
	// entries in this object to override those params from command line
	// CONFIG_FILE for this query. Most search parameters can be overriden:
	// cpuctExploration, winLossUtilityFactor, etc.
	//
	// priority (int): Optional. Analysis threads will prefer handling queries
	// with the highest priority unless already started on another task, breaking
	// ties in favor of earlier queries. If not specified, defaults to 0.
}

// New creates a Query object with default parameters
func NewQuery() *Query {
	return &Query{
		Id:    uuid.New().String(),
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
