package kataprob

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/otrego/clamshell/go/movetree"
	"github.com/otrego/clamshell/go/sgf"
	"github.com/otrego/clamshell/katago"
)

func TestBlunderComprehensiveness(t *testing.T) {
	game, err := MakeGameWithAnalysis(t)
	if err != nil {
		t.Fatal(err)
	}

	paths, err := FindBlunders(game)
	if err != nil {
		t.Fatal(err)
	}

	expected := 4
	actual := len(paths)
	if actual != expected {
		t.Errorf("Expected to find %d blunders, but found %d", expected, actual)
	}
}

func TestBlunderPointThreshold(t *testing.T) {
	game, err := MakeGameWithAnalysis(t)
	if err != nil {
		t.Fatal(err)
	}

	threshold := 8.0
	paths, err := FindBlundersWithOptions(game, &FindBlunderOptions{PointThreshold: threshold})
	if err != nil {
		t.Fatal(err)
	}

	expected := 3
	actual := len(paths)
	if actual != expected {
		t.Errorf("Expected to find %d blunders with a point delta > %f, but found %d", expected, threshold, actual)
	}
}

func MakeGameWithAnalysis(t *testing.T) (*movetree.MoveTree, error) {
	content, err := os.ReadFile("../../test-database/blunders.sgf")
	if err != nil {
		t.Fatal(err)
	}

	game, err := sgf.FromString(string(content)).Parse()
	if err != nil {
		t.Fatal(err)
	}

	analysisBytes, err := os.ReadFile("../testdata/blunder-analysis.json")
	if err != nil {
		t.Fatal(err)
	}

	analysis := &katago.AnalysisList{}
	if err := json.Unmarshal(analysisBytes, analysis); err != nil {
		t.Fatal(err)
	}

	if err := analysis.AddToGame(game); err != nil {
		t.Fatal(err)
	}

	return game, nil
}
