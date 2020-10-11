package katago

import (
	"encoding/json"
	"fmt"
	"math"
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
	BoardXSize int `json:"boardXSize,omitempty"`

	// BoardXSize is the width of the board
	BoardYSize int `json:"boardYSize,omitempty"`

	// AnalyzeTurns is the turns of the game to analyze. If this field is not
	// specified, defaults to analyzing the last turn only.
	AnalyzeTurns []int `json:"analyzeTurns,omitempty"`

	// MaxVisits is the maximum number of visits to use. If not specified,
	// defaults to the value in the analysis config file. If specified, overrides
	// it.
	MaxVisits *int `json:"maxVisits,omitempty"`

	OverrideSettings map[string]interface{} `json:"overrideSettings,omitempty"`

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
	// priority
}

// NewQuery creates an analysis Query object with default parameters
func NewQuery() *Query {
	return &Query{
		ID:               uuid.New().String(),
		Rules:            TrompTaylorRules,
		OverrideSettings: make(map[string]interface{}),
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

// gameConverter is used to convert games into query objects.
type gameConverter struct {
	g *game.Game
}

// point converts a point to a GTP point
func (gc *gameConverter) point(pt *point.Point) string {
	if pt == nil {
		// We treat nil points as passes.
		return "pass"
	}

	const a = 'A'
	const i = 'I'
	val := rune(a + pt.X())
	if val >= i {
		// I is strictly not allowed.
		val++
	}
	return fmt.Sprintf("%c%d", val, pt.Y()+1)
}

// move converts from a game-move to a move-array with a GTP Point. This is a
// format peculiar to Katago.
func (gc *gameConverter) move(mv *game.Move) Move {
	return Move{string(mv.Color()), gc.point(mv.Point())}
}

// initialStones converts the initial placements (e.g., handicap stones) into
// the initial stones for the analysis.
func (gc *gameConverter) initialStones() []Move {
	var out []Move
	for _, mv := range gc.g.Root.Placements {
		out = append(out, gc.move(mv))
	}
	return out
}

// initialPlayer sets the initial player. By default, katago assumes
// black-to-play, so this isn't necessary.
func (gc *gameConverter) initialPlayer() string {
	if val, ok := gc.g.Root.Properties["PL"]; ok && len(val) > 0 {
		return val[0]
	}
	return ""
}

// mainBranchMoves gets moves along the primary branch (0th-variation).
func (gc *gameConverter) mainBranchMoves() []Move {
	out := []Move{}
	idx := 0
	for n := gc.g.Root; ; n = n.Children[0] {
		if n.Move != nil {
			out = append(out, gc.move(n.Move))
		}
		if len(n.Children) == 0 {
			// No more children; terminate traversal.
			break
		}
		idx++
	}
	return out
}

// rules gets the relevant rule-set, returning TrompTaylorRules if not provided.
func (gc *gameConverter) rules() Rules {
	if val, ok := gc.g.Root.Properties["RU"]; ok && len(val) > 0 {
		return Rules(val[0])
	}
	return TrompTaylorRules
}

// komi gets the komi value, parsed as a float. Note that the decimal-part can
// only be exactl 0.5 or 0.
func (gc *gameConverter) komi() (*float64, error) {
	if val, ok := gc.g.Root.Properties["KM"]; ok && len(val) > 0 {
		km, err := strconv.ParseFloat(val[0], 64)
		if err != nil {
			return nil, err
		}
		_, fp := math.Modf(km)
		if !(fp == 0.5 || fp == 0.0) {
			return nil, fmt.Errorf("invalid komi: the only decimal-value allowed for komi is 0.0 or 0.5. komi was %f", km)
		}
		return &km, nil
	}
	return nil, nil
}

// boardSize gets the size of the go board. Only sizes ups to 25 are allowed,
// but should typically be 19, 13, or 9.
func (gc *gameConverter) boardSize() (int, error) {
	if val, ok := gc.g.Root.Properties["SZ"]; ok && len(val) > 0 {
		sz, err := strconv.Atoi(val[0])
		if err != nil {
			return 0, err
		}
		if sz > 25 {
			return 0, fmt.Errorf("only sizes up to 25 are supported")
		}
		return sz, nil
	}
	return 19, nil
}

// analyzeMainBranch analyzes the main branch of the game.
func (gc *gameConverter) analyzeMainBranch(startFrom, maxMoves int) []int {
	var out []int
	numAnalyzed := 0
	if maxMoves == 0 {
		maxMoves = 1000
	}
	for n := gc.g.Root; n != nil; n = n.Next(0) {
		if n.MoveNum() >= startFrom && numAnalyzed < maxMoves {
			out = append(out, n.MoveNum())
			numAnalyzed++
		}
	}
	return out
}

// NewInt creates an integer pointer.
func NewInt(i int) *int {
	return &i
}

// QueryOptions contains options for constructing a Query from a game.
type QueryOptions struct {
	// MaxMoves is used to provide a max when analyzing moves. By default, analyze the
	// whole game. Largely used for debugging.
	MaxMoves *int

	// StartFrom indicates which move to start from. By default, start from the
	// first move (but not the root).
	StartFrom *int

	// AnalysisDepth indicates how deep to print analyses. If not specified,
	// defaults to 0 (no depth to analysis).
	AnalysisDepth *int

	// MaxVisits indicates the maximum number of visits to perform during
	// analysis. If not specified, defaults to 10.
	MaxVisits *int
}

// defaultOptions adds defaults to query options
func defaultOptions(opts *QueryOptions) *QueryOptions {
	// defaults for the query options.
	d := &QueryOptions{
		MaxMoves:      NewInt(0),
		StartFrom:     NewInt(1),
		AnalysisDepth: NewInt(0),
		MaxVisits:     NewInt(10),
	}
	if opts == nil {
		opts = &QueryOptions{}
	}

	if opts.MaxMoves != nil {
		d.MaxMoves = opts.MaxMoves
	}
	if opts.StartFrom != nil {
		d.StartFrom = opts.StartFrom
	}
	if opts.AnalysisDepth != nil {
		d.AnalysisDepth = opts.AnalysisDepth
	}
	if opts.MaxVisits != nil {
		d.MaxVisits = opts.MaxVisits
	}

	return d
}

// AnalysisQueryFromGame creates a Katago query object from a go-game.
func AnalysisQueryFromGame(g *game.Game, inOpts *QueryOptions) (*Query, error) {
	q := NewQuery()
	gc := &gameConverter{g: g}
	opts := defaultOptions(inOpts)

	q.InitialStones = gc.initialStones()
	q.InitialPlayer = gc.initialPlayer()
	q.Moves = gc.mainBranchMoves()
	q.Rules = gc.rules()
	q.MaxVisits = opts.MaxVisits

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
	q.OverrideSettings["analysisPVLen"] = strconv.Itoa(*opts.AnalysisDepth)
	q.AnalyzeTurns = gc.analyzeMainBranch(*opts.StartFrom, *opts.MaxMoves)

	return q, nil
}
