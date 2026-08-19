package ring

import (
	"testing"

	"task116-chashring/internal/model"
)

func sampleNodes() []*model.Node {
	return []*model.Node{
		{ID: "a", Weight: 1, VirtualNodes: 10},
		{ID: "b", Weight: 1, VirtualNodes: 10},
		{ID: "c", Weight: 1, VirtualNodes: 10},
	}
}

func TestLookupDeterministic(t *testing.T) {
	rh := New("r", model.RingConfig{}, sampleNodes())
	p1, err := rh.Lookup("alpha")
	if err != nil {
		t.Fatal(err)
	}
	p2, err := rh.Lookup("alpha")
	if err != nil {
		t.Fatal(err)
	}
	if p1 != p2 {
		t.Fatalf("lookup not deterministic: %s != %s", p1, p2)
	}
	if p1 == "" {
		t.Fatal("empty lookup result")
	}
}

func TestAddRemoveNode(t *testing.T) {
	rh := New("r", model.RingConfig{}, sampleNodes())
	if err := rh.AddNode(&model.Node{ID: "d", Weight: 1, VirtualNodes: 10}); err != nil {
		t.Fatal(err)
	}
	if len(rh.Nodes()) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(rh.Nodes()))
	}
	if err := rh.RemoveNode("d"); err != nil {
		t.Fatal(err)
	}
	if len(rh.Nodes()) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(rh.Nodes()))
	}
	if err := rh.RemoveNode("missing"); err == nil {
		t.Fatal("expected not-found error")
	}
}

func TestStatsCounts(t *testing.T) {
	rh := New("r", model.RingConfig{}, sampleNodes())
	s := rh.Stats()
	if s.NodeCount != 3 {
		t.Fatalf("expected 3 nodes, got %d", s.NodeCount)
	}
	if s.PointCount == 0 {
		t.Fatal("expected non-zero point count")
	}
}

func TestReplicasSingle(t *testing.T) {
	rh := New("r", model.RingConfig{}, sampleNodes())
	reps, err := rh.Replicas("alpha", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(reps) != 1 {
		t.Fatalf("expected 1 replica, got %d", len(reps))
	}
}

func TestReplicasAllDistinct(t *testing.T) {
	rh := New("r", model.RingConfig{}, sampleNodes())
	reps, err := rh.Replicas("alpha", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(reps) != 3 {
		t.Fatalf("expected 3 distinct replicas, got %d", len(reps))
	}
	seen := map[string]bool{}
	for _, r := range reps {
		if seen[r] {
			t.Fatalf("duplicate replica %s", r)
		}
		seen[r] = true
	}
}
