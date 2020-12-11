package katago

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/otrego/clamshell/core/movetree"
	"github.com/otrego/clamshell/core/sgf"
)

func TestParseResult_Short(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/short-analysis.json")
	if err != nil {
		t.Fatal(err)
	}
	resList, err := ParseAnalysisList(content)
	if err != nil {
		t.Fatal(err)
	}
	if len(resList) != 1 {
		t.Errorf("Got %d documents, expected exactly 1", len(resList))
	}

	item := resList[0]
	expID := "863f6928-6b84-45e7-989e-af847c015320"
	if item.ID != expID {
		t.Errorf("Got ID %s, but expected %s", item.ID, expID)
	}

	expTurnNumber := 10
	if item.TurnNumber != expTurnNumber {
		t.Errorf("Got turn number %d, but expected %d", item.TurnNumber, expTurnNumber)
	}

	if item.MoveInfos == nil {
		t.Error("MoveInfos was nil, but was expected to have values")
	}

	if item.RootInfo == nil {
		t.Error("RootInfo was nil, but was expected to have values")
	}
}

func TestParseResult_Long(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/long-analysis.json")
	if err != nil {
		t.Fatal(err)
	}
	resList, err := ParseAnalysisList(content)
	if err != nil {
		t.Fatal(err)
	}
	if len(resList) != 240 {
		t.Errorf("Got %d documents, expected exactly 240", len(resList))
	}

	item := resList[0]
	expID := "93f7981f-0b5b-423f-90ac-f7961622b9bd"
	if item.ID != expID {
		t.Errorf("Got ID %s, but expected %s", item.ID, expID)
	}

	expTurnNumber := 1
	if item.TurnNumber != expTurnNumber {
		t.Errorf("Got turn number %d, but expected %d", item.TurnNumber, expTurnNumber)
	}

	if item.RootInfo == nil {
		t.Error("RootInfo was nil, but was expected to have values")
	}
}

func TestDemarshalMoveInfo_Single(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/moveinfo_single.json")
	if err != nil {
		t.Fatal(err)
	}
	dec := json.NewDecoder(strings.NewReader(string(content)))
	res := &MoveInfo{}
	err = dec.Decode(res)
	if err != nil {
		t.Fatal(err)
	}
}
func Test_Demarshal_MoveInfo_Single(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/moveinfo_single.json")
	if err != nil {
		t.Fatal(err)
	}
	dec := json.NewDecoder(strings.NewReader(string(content)))
	res := &MoveInfo{}
	err = dec.Decode(res)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Demarshal_MoveInfo_Slice(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/moveinfo_multi.json")
	if err != nil {
		t.Fatal(err)
	}
	dec := json.NewDecoder(strings.NewReader(string(content)))
	res := []MoveInfo{}
	err = dec.Decode(&res)
	if err != nil {
		t.Fatal(err)
	}
	expectedMoves := 16
	if len(res) != expectedMoves {
		t.Errorf("Got %d moves, but expected %d", len(res), expectedMoves)
	}
}

func TestSortResults(t *testing.T) {
	res := AnalysisList{
		&AnalysisResult{ID: "foo-3", TurnNumber: 3},
		&AnalysisResult{ID: "foo-5", TurnNumber: 5},
		&AnalysisResult{ID: "foo-1", TurnNumber: 1},
		&AnalysisResult{ID: "foo-4", TurnNumber: 4},
	}
	res.Sort()
	exp := []int{1, 3, 4, 5}

	var got []int
	for _, e := range res {
		got = append(got, e.TurnNumber)
	}
	if !cmp.Equal(got, exp) {
		t.Errorf("got %v, but expected %v; diff %v", got, exp, cmp.Diff(got, exp))
	}
}

func TestAddToGame(t *testing.T) {
	floatPtr := func(f float64) *float64 {
		return &f
	}

	testCases := []struct {
		desc     string
		rawSgf   string
		analysis AnalysisList

		// map of treepath string to expected winrate value
		expWinRate map[string]*float64
	}{
		{
			desc:   "basic attachment",
			rawSgf: `(;GM[1];B[ac];W[cd];B[de])`,
			analysis: AnalysisList{
				&AnalysisResult{
					ID:         "foo",
					TurnNumber: 1,
					RootInfo: &RootInfo{
						Winrate: 0.1,
					},
				},
				&AnalysisResult{
					ID:         "foo",
					TurnNumber: 3,
					RootInfo: &RootInfo{
						Winrate: 0.3,
					},
				},
			},
			expWinRate: map[string]*float64{
				"-":      nil,
				"-0":     floatPtr(0.1),
				"-0-0":   nil,
				"-0-0-0": floatPtr(0.3),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := sgf.Parse(tc.rawSgf)
			if err != nil {
				t.Fatal(err)
			}
			if err := tc.analysis.AddToGame(g); err != nil {
				t.Fatal(err)
			}
			for tpRaw, valp := range tc.expWinRate {
				tp, err := movetree.ParsePath(tpRaw)
				if err != nil {
					t.Error(err)
					continue
				}
				n := tp.Apply(g.Root)

				ad := n.AnalysisData()
				if valp == nil && ad == nil {
					// expected case, but nothing to do.
					continue
				} else if valp == nil && ad != nil {
					t.Errorf("at treepath %q, got analysis data, but expected none", tpRaw)
					continue
				} else if valp != nil && ad == nil {
					t.Errorf("at treepath %q, got no analysis data, but expected some", tpRaw)
					continue
				}
				val := *valp

				nodeAn, ok := n.AnalysisData().(*AnalysisResult)
				if !ok {
					t.Errorf("at treepath %q, attached analysis was of wrong type. Type was %T", tpRaw, nodeAn)
					continue
				}
				if wr := nodeAn.RootInfo.Winrate; wr != val {
					t.Errorf("at treepath %q, got winrate %f, but expected %f", tpRaw, wr, val)
				}
			}
		})
	}
}
