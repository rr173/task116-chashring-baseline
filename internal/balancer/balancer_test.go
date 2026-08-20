package balancer

import (
	"fmt"
	"testing"

	"task116-chashring/internal/model"
)

func TestPreviewNoChange(t *testing.T) {
	cfg := model.RingConfig{Replicas: 20, Replication: 2}
	current := []model.Node{
		{ID: "a", Weight: 1, VirtualNodes: 20},
		{ID: "b", Weight: 1, VirtualNodes: 20},
	}
	sample := []string{"k1", "k2", "k3"}
	plan := Preview(cfg, current, sample, nil, nil)
	if plan.Count != 0 {
		t.Fatalf("expected no moves, got %d", plan.Count)
	}
}

func TestPreviewWithAdd(t *testing.T) {
	cfg := model.RingConfig{Replicas: 20, Replication: 2}
	current := []model.Node{{ID: "a", Weight: 1, VirtualNodes: 20}}
	sample := make([]string, 0, 200)
	for i := 0; i < 200; i++ {
		sample = append(sample, fmt.Sprintf("k%d", i))
	}
	added := []model.Node{{ID: "b", Weight: 1, VirtualNodes: 20}}
	plan := Preview(cfg, current, sample, added, nil)
	if plan.Count == 0 {
		t.Fatalf("expected some moves after adding a node")
	}
	if len(plan.Added) != 1 {
		t.Fatalf("added not recorded")
	}
}

// TestPreviewDeduplicatesKeys verifies that duplicate, whitespace-padded, and
// blank sample keys collapse to a single move per meaningful key. Node "a" is
// removed while "b" is added, so every surviving key moves a→b deterministically.
func TestPreviewDeduplicatesKeys(t *testing.T) {
	cfg := model.RingConfig{Replicas: 20, Replication: 2}
	current := []model.Node{{ID: "a", Weight: 1, VirtualNodes: 20}}
	added := []model.Node{{ID: "b", Weight: 1, VirtualNodes: 20}}
	removed := []string{"a"}
	sample := []string{"alpha", " alpha ", "", "alpha", "\t"}
	plan := Preview(cfg, current, sample, added, removed)
	if plan.Count != 1 {
		t.Fatalf("expected one move, got %d (%#v)", plan.Count, plan.Moves)
	}
	if plan.Moves[0].Key != "alpha" {
		t.Fatalf("expected move for alpha, got %q", plan.Moves[0].Key)
	}
	if plan.Moves[0].From != "a" || plan.Moves[0].To != "b" {
		t.Fatalf("unexpected move: %+v", plan.Moves[0])
	}
}
