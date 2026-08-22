package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"task116-chashring/internal/model"
	"task116-chashring/internal/store"
)

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func newTestAPI(t *testing.T) *API {
	t.Helper()
	a := New(newTestStore(t))
	if err := a.LoadFromStore(context.Background()); err != nil {
		t.Fatal(err)
	}
	return a
}

func doJSON(t *testing.T, srv *httptest.Server, method, path string, body any) *http.Response {
	t.Helper()
	var r *http.Request
	var err error
	if body != nil {
		b, _ := json.Marshal(body)
		r, err = http.NewRequest(method, srv.URL+path, bytes.NewReader(b))
	} else {
		r, err = http.NewRequest(method, srv.URL+path, nil)
	}
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func TestCreateRingAndNode(t *testing.T) {
	a := newTestAPI(t)
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()

	resp := doJSON(t, srv, "POST", "/rings", model.Ring{ID: "r1", Name: "ring", Config: model.RingConfig{Replicas: 20}})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create status %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp = doJSON(t, srv, "POST", "/rings/r1/nodes", model.Node{ID: "n1", Address: "x:1", Weight: 1, VirtualNodes: 20})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("add node status %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp = doJSON(t, srv, "GET", "/rings/r1", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get ring status %d", resp.StatusCode)
	}
	var view model.RingView
	if err := json.NewDecoder(resp.Body).Decode(&view); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if len(view.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(view.Nodes))
	}

	resp = doJSON(t, srv, "GET", "/rings", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list rings status %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp = doJSON(t, srv, "GET", "/rings/r1/lookup?key=alpha", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("lookup status %d", resp.StatusCode)
	}
	resp.Body.Close()

	resp = doJSON(t, srv, "GET", "/rings/r1/stats", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("stats status %d", resp.StatusCode)
	}
	resp.Body.Close()
}

// TestRemoveNodeUpdatesBothStates asserts that deleting a node removes it from
// the persisted store and from the live routing ring in lockstep. This guards
// against a regression where the storage row was dropped but the node kept
// routing traffic in the in-memory ring.
func TestRemoveNodeUpdatesBothStates(t *testing.T) {
	a := newTestAPI(t)
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()

	if resp := doJSON(t, srv, "POST", "/rings", model.Ring{ID: "r2", Name: "r2", Config: model.RingConfig{Replicas: 20, Replication: 3}}); resp.StatusCode != http.StatusCreated {
		t.Fatalf("create ring status %d", resp.StatusCode)
	} else {
		resp.Body.Close()
	}
	for _, id := range []string{"n1", "n2"} {
		resp := doJSON(t, srv, "POST", "/rings/r2/nodes", model.Node{ID: id, Address: "x:1", Weight: 1, VirtualNodes: 20})
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("add node %s status %d", id, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// Both nodes must be live and persisted before the delete.
	rh, _ := a.Ring("r2")
	if got := len(rh.Nodes()); got != 2 {
		t.Fatalf("expected 2 live nodes before delete, got %d", got)
	}

	resp := doJSON(t, srv, "DELETE", "/rings/r2/nodes/n1", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete status %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Storage state must no longer contain n1.
	persisted, err := a.Store().ListNodes(context.Background(), "r2")
	if err != nil {
		t.Fatal(err)
	}
	if len(persisted) != 1 || persisted[0].ID != "n2" {
		t.Fatalf("expected only n2 persisted, got %+v", persisted)
	}

	// Routing ring state must no longer contain n1.
	live := rh.Nodes()
	if len(live) != 1 || live[0].ID != "n2" {
		t.Fatalf("expected only n2 live, got %+v", live)
	}

	// A lookup must never resolve to the deleted node.
	node, err := rh.Lookup("user-session-abc")
	if err != nil {
		t.Fatalf("lookup after delete: %v", err)
	}
	if node == "n1" {
		t.Fatal("lookup resolved to deleted node n1")
	}

	// Re-deleting must surface not-found (storage is the source of truth).
	if err := a.RemoveNode(context.Background(), "r2", "n1"); err != model.ErrNodeNotFound {
		t.Fatalf("expected ErrNodeNotFound on second delete, got %v", err)
	}
}
