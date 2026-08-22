package store

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"task116-chashring/internal/model"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	st, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func TestRingCRUD(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	r := model.Ring{ID: "r1", Name: "ring", Config: model.RingConfig{Replicas: 10, Replication: 2}, CreatedAt: 1, UpdatedAt: 1}
	if err := st.SaveRing(ctx, r); err != nil {
		t.Fatal(err)
	}
	got, err := st.GetRing(ctx, "r1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "ring" {
		t.Fatalf("name mismatch: %s", got.Name)
	}
	if _, err := st.GetRing(ctx, "nope"); err == nil {
		t.Fatal("expected not found")
	}
	rings, err := st.ListRings(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rings) != 1 {
		t.Fatalf("expected 1 ring, got %d", len(rings))
	}
}

func TestSummarySeparatesNodeAndPointCounts(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	r := model.Ring{ID: "r1", Config: model.RingConfig{Replicas: 10, Replication: 2}, CreatedAt: 1, UpdatedAt: 1}
	if err := st.SaveRing(ctx, r); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		n := model.Node{ID: fmt.Sprintf("n%d", i), Address: "x", Weight: 1, VirtualNodes: 10, CreatedAt: int64(i)}
		if err := st.SaveNode(ctx, "r1", n); err != nil {
			t.Fatal(err)
		}
	}
	summary, err := st.Summary(ctx, "r1")
	if err != nil {
		t.Fatal(err)
	}
	// Nodes is the number of distinct physical backends.
	if summary.Nodes != 3 {
		t.Fatalf("expected 3 physical nodes, got %d", summary.Nodes)
	}
	// Points is the total number of virtual routing points: 3 nodes * 10 each.
	if summary.Points != 30 {
		t.Fatalf("expected 30 virtual points, got %d", summary.Points)
	}
	if summary.Nodes == summary.Points {
		t.Fatal("summary collapses physical node and virtual point counts into one number")
	}
}

func TestNodeDuplicate(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	r := model.Ring{ID: "r1", Config: model.RingConfig{}, CreatedAt: 1, UpdatedAt: 1}
	if err := st.SaveRing(ctx, r); err != nil {
		t.Fatal(err)
	}
	n := model.Node{ID: "n1", Address: "x", Weight: 1, VirtualNodes: 5, CreatedAt: 1}
	if err := st.SaveNode(ctx, "r1", n); err != nil {
		t.Fatal(err)
	}
	if err := st.SaveNode(ctx, "r1", n); err != model.ErrNodeDuplicate {
		t.Fatalf("expected duplicate error, got %v", err)
	}
	nodes, err := st.ListNodes(ctx, "r1")
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	if err := st.DeleteNode(ctx, "r1", "n1"); err != nil {
		t.Fatal(err)
	}
	if err := st.DeleteNode(ctx, "r1", "n1"); err != model.ErrNodeNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}
