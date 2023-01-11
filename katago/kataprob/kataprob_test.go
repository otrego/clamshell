package kataprob

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/otrego/clamshell/go/sgf"
	"github.com/otrego/clamshell/katago"
)

func TestBlunderComprehensiveness(t *testing.T) {
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

	paths, err := FindBlunders(game)
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 4 {
		t.Error("Expected to find 4 blunders")
	}
}
