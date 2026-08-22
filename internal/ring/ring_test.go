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
	// NodeCount is the number of distinct physical nodes.
	if s.NodeCount != 3 {
		t.Fatalf("expected 3 physical nodes, got %d", s.NodeCount)
	}
	// PointCount is the total number of virtual routing points: 3 nodes * 10
	// virtual nodes each, weight 1.
	if s.PointCount != 30 {
		t.Fatalf("expected 30 virtual points, got %d", s.PointCount)
	}
	// Each physical node must report its own virtual-point contribution.
	totalVN := 0
	for _, nd := range s.Nodes {
		if nd.VirtualNodes != 10 {
			t.Fatalf("node %s: expected 10 virtual points, got %d", nd.ID, nd.VirtualNodes)
		}
		totalVN += nd.VirtualNodes
	}
	if totalVN != s.PointCount {
		t.Fatalf("per-node virtual points %d must sum to PointCount %d", totalVN, s.PointCount)
	}
}

func TestMetricsCounts(t *testing.T) {
	rh := New("r", model.RingConfig{}, sampleNodes())
	m := rh.Metrics()
	if m.Nodes != 3 {
		t.Fatalf("expected 3 physical nodes, got %d", m.Nodes)
	}
	if m.Points != 30 {
		t.Fatalf("expected 30 virtual points, got %d", m.Points)
	}
	if m.Nodes == m.Points {
		t.Fatal("metrics collapse node and point counts into one number")
	}
}

func TestHealthCounts(t *testing.T) {
	rh := New("r", model.RingConfig{}, sampleNodes())
	h := rh.Health()
	if h.Nodes != 3 {
		t.Fatalf("expected 3 physical nodes, got %d", h.Nodes)
	}
	if h.Points != 30 {
		t.Fatalf("expected 30 virtual points, got %d", h.Points)
	}
	if h.Empty {
		t.Fatal("non-empty ring reported as empty")
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
