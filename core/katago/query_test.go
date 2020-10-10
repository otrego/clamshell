package katago

import (
	"testing"
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
