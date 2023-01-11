package kataprob

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/otrego/clamshell/go/color"
	"github.com/otrego/clamshell/go/movetree"
	"github.com/otrego/clamshell/go/sgf"
	"github.com/otrego/clamshell/katago"
)

func TestBlunderComprehensiveness(t *testing.T) {
	game, err := makeGameWithAnalysis(t)
	if err != nil {
		t.Fatal(err)
	}

	paths, err := FindBlunders(game)
	if err != nil {
		t.Fatal(err)
	}

	expectedBlunders := []int{2, 3, 5, 6}
	for i, expected := range expectedBlunders {
		actual := paths[i].Apply(game.Root).MoveNum()
		if actual != expected {
			t.Errorf("Expected to find blunder at move %d, but found move %d", expected, actual)
		}
	}
}

func TestBlunderPointThreshold(t *testing.T) {
	game, err := makeGameWithAnalysis(t)
	if err != nil {
		t.Fatal(err)
	}

	threshold := 8.0
	paths, err := FindBlundersWithOptions(game, &FindBlunderOptions{PointThreshold: threshold})
	if err != nil {
		t.Fatal(err)
	}

	expectedBlunders := []int{2, 3, 6}
	for i, expected := range expectedBlunders {
		actual := paths[i].Apply(game.Root).MoveNum()
		if actual != expected {
			t.Errorf("Expected to find blunder at move %d, but found move %d", expected, actual)
		}
	}
}

func TestBlunderBlackColorFilter(t *testing.T) {
	expectedColor := color.Black
	game, err := makeGameWithAnalysis(t)
	if err != nil {
		t.Fatal(err)
	}

	threshold := 3.0
	paths, err := FindBlundersWithOptions(
		game,
		&FindBlunderOptions{PointThreshold: threshold, Color: expectedColor},
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedBlunders := []int{3, 5}
	for i, expected := range expectedBlunders {
		root := paths[i].Apply(game.Root)
		actual := root.MoveNum()
		if actual != expected {
			t.Errorf("Expected to find blunder at move %d, but found move %d", expected, actual)
		}
		actualColor := root.Move.Color()
		if actualColor != expectedColor {
			t.Errorf("Found color %#v for move %d, expected %#v", actualColor, root.MoveNum(), expectedColor)
		}
	}
}

func TestBlunderWhiteColorFilter(t *testing.T) {
	expectedColor := color.White
	game, err := makeGameWithAnalysis(t)
	if err != nil {
		t.Fatal(err)
	}

	threshold := 3.0
	paths, err := FindBlundersWithOptions(
		game,
		&FindBlunderOptions{PointThreshold: threshold, Color: expectedColor},
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedBlunders := []int{2, 6}
	for i, expected := range expectedBlunders {
		root := paths[i].Apply(game.Root)
		actual := root.MoveNum()
		if actual != expected {
			t.Errorf("Expected to find blunder at move %d, but found move %d", expected, actual)
		}
		actualColor := root.Move.Color()
		if actualColor != expectedColor {
			t.Errorf("Found color %#v for move %d, expected %#v", actualColor, root.MoveNum(), expectedColor)
		}
	}
}

func makeGameWithAnalysis(t *testing.T) (*movetree.MoveTree, error) {
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
