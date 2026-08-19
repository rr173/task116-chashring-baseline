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
