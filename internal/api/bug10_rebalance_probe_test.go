package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"task116-chashring/internal/balancer"
	"task116-chashring/internal/model"
)

func TestBug10_RebalancePreviewAndApplyBothRemoveSelectedNode(t *testing.T) {
	ctx := context.Background()
	a := newTestAPI(t)
	if err := a.CreateRing(ctx, model.Ring{ID: "r", Name: "r", Config: model.RingConfig{Replicas: 7}}); err != nil { t.Fatal(err) }
	for _, id := range []string{"a", "b"} { if err := a.AddNode(ctx, "r", model.Node{ID: id, Address: id + ":1", Weight: 1}); err != nil { t.Fatal(err) } }
	rh, _ := a.Ring("r")
	removed, err := rh.Lookup("alpha")
	if err != nil { t.Fatal(err) }
	plan := balancer.Preview(rh.Config(), rh.Nodes(), []string{"alpha"}, nil, []string{removed})
	if plan.Count != 1 { t.Fatalf("preview count = %d, want 1", plan.Count) }
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()
	resp := doJSON(t, srv, http.MethodPost, "/rings/r/rebalance", map[string]any{"removed": []string{removed}})
	if resp.StatusCode != http.StatusOK { t.Fatalf("rebalance status = %d", resp.StatusCode) }
	resp.Body.Close()
	if rh.HasNode(removed) { t.Fatalf("removed node %q is still live", removed) }
}
