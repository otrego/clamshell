package katago

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/otrego/clamshell/core/errcheck"
	"github.com/otrego/clamshell/core/sgf"
)

func TestCreateQuery_Defaults(t *testing.T) {
	q := NewQuery()
	if q == nil {
		t.Fatal("nil query")
	}
	if q.ID == "" {
		t.Error("empty ID for query")
	}
	if q.Rules != TrompTaylorRules {
		t.Errorf("got rules %v, but expected %v", q.Rules, TrompTaylorRules)
	}

	_, err := q.ToJSON()
	if err != nil {
		t.Errorf("unexpected error during marshaling: %v", err)
	}
}

func TestCreateAnalysis(t *testing.T) {
	defaultQuery := func() *Query {
		return &Query{
			ID:        "fakeid",
			Rules:     TrompTaylorRules,
			MaxVisits: NewInt(10),
			Moves:     []Move{},
			OverrideSettings: map[string]interface{}{
				"analysisPVLen": "0",
			},
			BoardXSize: 19,
			BoardYSize: 19,
		}
	}

	testCases := []struct {
		desc         string
		sgf          string
		opts         *QueryOptions
		expQuery     *Query
		expErrSubstr string
	}{
		{
			desc:     "basic case",
			sgf:      "(;GM[1])",
			expQuery: defaultQuery(),
		},
		{
			desc: "basic case: analyze root",
			sgf:  "(;GM[1])",
			opts: &QueryOptions{
				StartFrom: NewInt(0),
			},
			expQuery: func() *Query {
				q := defaultQuery()
				q.AnalyzeTurns = []int{0}
				return q
			}(),
		},
		{
			desc: "Analyze some moves",
			sgf:  "(;GM[1];B[aa];W[bb];B[cc];W[dd])",
			expQuery: func() *Query {
				q := defaultQuery()
				q.Moves = []Move{
					Move{"B", "A1"},
					Move{"W", "B2"},
					Move{"B", "C3"},
					Move{"W", "D4"},
				}
				q.AnalyzeTurns = []int{1, 2, 3, 4}
				return q
			}(),
		},
		{
			desc: "Analyze some moves: Start From",
			sgf:  "(;GM[1];B[aa];W[bb];B[cc];W[dd])",
			opts: &QueryOptions{
				StartFrom: NewInt(3),
			},
			expQuery: func() *Query {
				q := defaultQuery()
				q.Moves = []Move{
					Move{"B", "A1"},
					Move{"W", "B2"},
					Move{"B", "C3"},
					Move{"W", "D4"},
				}
				q.AnalyzeTurns = []int{3, 4}
				return q
			}(),
		},
		{
			desc: "Analyze some moves: Max moves",
			sgf:  "(;GM[1];B[aa];W[bb];B[cc];W[dd])",
			opts: &QueryOptions{
				MaxMoves: NewInt(2),
			},
			expQuery: func() *Query {
				q := defaultQuery()
				q.Moves = []Move{
					Move{"B", "A1"},
					Move{"W", "B2"},
				}
				q.AnalyzeTurns = []int{1, 2}
				return q
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := sgf.Parse(tc.sgf)
			if err != nil {
				t.Fatal(err)
			}
			got, err := AnalysisQueryFromGame(g, tc.opts)
			cerr := errcheck.CheckCases(err, tc.expErrSubstr)
			if cerr != nil {
				t.Fatal(cerr)
			}
			if err != nil {
				return
			}

			// Set some values to improve the equals/diff
			got.ID = "fakeid"

			gj, err := got.ToJSON()
			if err != nil {
				t.Fatal(err)
			}
			ej, err := tc.expQuery.ToJSON()
			if err != nil {
				t.Fatal(err)
			}
			gjs, ejs := string(gj), string(ej)
			if !cmp.Equal(gjs, ejs) {
				t.Errorf("AnalysisQueryFromGame =>\n%v\nwanted\n%v\ndiff=%v", gjs, ejs, cmp.Diff(ejs, gjs))
			}

			// For pasting into katago:
			// t.Error(gjs)
		})
	}
}
