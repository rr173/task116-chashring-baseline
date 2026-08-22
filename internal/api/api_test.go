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

// TestDefaultConfigPersistsAcrossReload reproduces the reported bug: creating a
// ring with the default (empty) config left replicas=0 on disk and built the
// in-memory ring with Replicas=1, so a restart recovered a different ring and
// /config reported the wrong capacity. After the fix the effective defaults
// (Replicas=100, Replication=2) must persist, /config must report them, and the
// ring rebuilt from the store must be byte-for-byte identical to the original.
func TestDefaultConfigPersistsAcrossReload(t *testing.T) {
	a := newTestAPI(t)
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()

	// Create a ring with the default/empty config.
	resp := doJSON(t, srv, "POST", "/rings", model.Ring{ID: "def", Name: "defaults"})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create status %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Add a couple of nodes whose own VirtualNodes are unset so they inherit
	// the ring's Replicas — this is exactly the path the old bug corrupted.
	for _, n := range []model.Node{
		{ID: "n0", Address: "10.0.0.0:8080", Weight: 1},
		{ID: "n1", Address: "10.0.0.1:8080", Weight: 1},
		{ID: "n2", Address: "10.0.0.2:8080", Weight: 1},
	} {
		resp := doJSON(t, srv, "POST", "/rings/def/nodes", n)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("add node %s status %d", n.ID, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// In-memory ring reflects the effective defaults, not Replicas=1.
	rh, ok := a.Ring("def")
	if !ok {
		t.Fatal("ring missing after create")
	}
	memCfg := rh.Config()
	if memCfg.Replicas != 100 || memCfg.Replication != 2 || memCfg.HashFunc != "fnv1a" {
		t.Fatalf("in-memory config not defaulted: %+v", memCfg)
	}

	// /config reads from the store and must report the same effective defaults.
	resp = doJSON(t, srv, "GET", "/rings/def/config", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("config status %d", resp.StatusCode)
	}
	var cfgResp struct {
		Config model.RingConfig `json:"config"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&cfgResp); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if cfgResp.Config.Replicas != 100 {
		t.Fatalf("expected /config replicas=100, got %d", cfgResp.Config.Replicas)
	}

	// Simulate a restart: a fresh API reloads from the same store.
	a2 := New(a.store)
	if err := a2.LoadFromStore(context.Background()); err != nil {
		t.Fatal(err)
	}
	rh2, ok := a2.Ring("def")
	if !ok {
		t.Fatal("ring missing after reload")
	}
	// The same lookup key must land on the same primary before and after
	// reload — this only holds if the rebuilt ring uses the identical config.
	p1, err := rh.Lookup("some-key")
	if err != nil {
		t.Fatal(err)
	}
	p2, err := rh2.Lookup("some-key")
	if err != nil {
		t.Fatal(err)
	}
	if p1 != p2 {
		t.Fatalf("lookup diverged after reload: %s != %s", p1, p2)
	}
	if rh2.Config() != memCfg {
		t.Fatalf("config changed after reload: %+v != %+v", rh2.Config(), memCfg)
	}
}
