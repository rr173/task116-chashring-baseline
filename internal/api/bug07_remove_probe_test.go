package api

import (
	"context"
	"testing"

	"task116-chashring/internal/model"
)

func TestBug07_NodeRemovalLeavesNoPersistentOrLiveOwner(t *testing.T) {
	ctx := context.Background()
	a := newTestAPI(t)
	if err := a.CreateRing(ctx, model.Ring{ID: "r", Name: "r"}); err != nil { t.Fatal(err) }
	if err := a.AddNode(ctx, "r", model.Node{ID: "n", Address: "n:1", Weight: 1}); err != nil { t.Fatal(err) }
	if err := a.RemoveNode(ctx, "r", "n"); err != nil { t.Fatal(err) }
	if _, err := a.Store().GetNode(ctx, "r", "n"); err != model.ErrNodeNotFound { t.Fatalf("stored node error = %v, want not found", err) }
	if rh, _ := a.Ring("r"); rh.HasNode("n") { t.Fatal("removed node is still live") }
}
