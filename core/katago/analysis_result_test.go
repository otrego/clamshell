package katago

import (
	"io/ioutil"
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
}
