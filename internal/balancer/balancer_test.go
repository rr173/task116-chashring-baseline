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
