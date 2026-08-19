package api

import (
	"context"
	"testing"

	"task116-chashring/internal/model"
)

func TestBug08_ReloadRestoresPersistedNodesIntoLiveRing(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	if err := st.SaveRing(ctx, model.Ring{ID: "r", Name: "r", Config: model.RingConfig{Replicas: 5}, CreatedAt: 1, UpdatedAt: 1}); err != nil { t.Fatal(err) }
	if err := st.SaveNode(ctx, "r", model.Node{ID: "n", Address: "n:1", Weight: 1, CreatedAt: 1}); err != nil { t.Fatal(err) }
	a := New(st)
	if err := a.LoadFromStore(ctx); err != nil { t.Fatal(err) }
	rh, ok := a.Ring("r")
	if !ok || !rh.HasNode("n") { t.Fatal("reload did not restore the persisted node") }
}
