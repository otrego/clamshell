package katago

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
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

	if item.MoveInfo == nil {
		t.Error("MoveInfo was nil, but was expected to have values")
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
