// Package selfcheck provides an offline, dependency-free end-to-end smoke test.
package selfcheck

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"task116-chashring/internal/api"
	"task116-chashring/internal/model"
	"task116-chashring/internal/store"
)

// Run builds a ring, exercises core operations, and verifies that state survives
// a store reload (simulating a process restart). It returns nil on success.
func Run(ctx context.Context) error {
	dir, err := os.MkdirTemp("", "chashring-smoke-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	db := filepath.Join(dir, "smoke.db")

	st, err := store.Open(db)
	if err != nil {
		return err
	}
	defer st.Close()

	a := api.New(st)
	r := model.Ring{ID: "smoke", Name: "smoke", Config: model.RingConfig{Replicas: 20, Replication: 3}}
	if err := a.CreateRing(ctx, r); err != nil {
		return fmt.Errorf("create ring: %w", err)
	}
	for i := 0; i < 5; i++ {
		n := model.Node{
			ID:           fmt.Sprintf("n%d", i),
			Address:      fmt.Sprintf("10.0.0.%d:8080", i),
			Weight:       1,
			VirtualNodes: 20,
		}
		if err := a.AddNode(ctx, "smoke", n); err != nil {
			return fmt.Errorf("add node: %w", err)
		}
	}

	rh, ok := a.Ring("smoke")
	if !ok {
		return fmt.Errorf("ring missing after create")
	}

	key := "user-session-abc"
	p1, err := rh.Lookup(key)
	if err != nil {
		return fmt.Errorf("lookup: %w", err)
	}
	p2, err := rh.Lookup(key)
	if err != nil {
		return fmt.Errorf("lookup2: %w", err)
	}
	if p1 != p2 {
		return fmt.Errorf("lookup not stable: %s != %s", p1, p2)
	}

	reps, err := rh.Replicas(key, 3)
	if err != nil {
		return fmt.Errorf("replicas: %w", err)
	}
	if len(reps) != 3 {
		return fmt.Errorf("expected 3 replicas, got %d", len(reps))
	}

	// Persistence + restart recovery: a fresh API loads the same state.
	a2 := api.New(st)
	if err := a2.LoadFromStore(ctx); err != nil {
		return fmt.Errorf("reload: %w", err)
	}
	rh2, ok := a2.Ring("smoke")
	if !ok {
		return fmt.Errorf("ring missing after reload")
	}
	p3, err := rh2.Lookup(key)
	if err != nil {
		return fmt.Errorf("lookup after reload: %w", err)
	}
	if p3 != p1 {
		return fmt.Errorf("lookup changed after reload: %s != %s", p3, p1)
	}
	return nil
}
