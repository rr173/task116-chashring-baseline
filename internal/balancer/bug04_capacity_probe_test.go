package balancer

import (
	"testing"

	"task116-chashring/internal/model"
)

func TestBug04_DefaultNodeCapacityMatchesBalanceTotals(t *testing.T) {
	nodes := []model.Node{{ID: "a", Weight: 2}, {ID: "b", Weight: 1}}
	if got := model.Capacity(nodes[0]); got != 200 { t.Fatalf("node capacity = %d, want 200", got) }
	if got := Weights(nodes).Total; got != 300 { t.Fatalf("weight report total = %d, want 300", got) }
}
