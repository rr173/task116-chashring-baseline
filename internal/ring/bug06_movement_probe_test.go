package ring

import (
	"testing"

	"task116-chashring/internal/model"
)

func TestBug06_RebalanceReportsEachMeaningfulKeyOnce(t *testing.T) {
	before := New("before", model.RingConfig{Replicas: 7}, []*model.Node{{ID: "old", Address: "old:1"}})
	after := New("after", model.RingConfig{Replicas: 7}, []*model.Node{{ID: "new", Address: "new:1"}})
	moves := before.Compare([]string{"alpha", " alpha ", "", "alpha"}, after)
	if len(moves) != 1 || moves[0].Key != "alpha" {
		t.Fatalf("moves = %#v, want one move for alpha", moves)
	}
}
