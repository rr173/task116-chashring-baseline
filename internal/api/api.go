// Package api exposes the consistent-hash ring service over HTTP.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"task116-chashring/internal/balancer"
	"task116-chashring/internal/model"
	"task116-chashring/internal/ring"
	"task116-chashring/internal/store"
)

// API holds the in-memory ring index and the persistence store.
type API struct {
	mu    sync.RWMutex
	store *store.Store
	rings map[string]*ring.ConsistentHash
}

// New constructs an API backed by the given store.
func New(s *store.Store) *API {
	return &API{store: s, rings: make(map[string]*ring.ConsistentHash)}
}

// Store returns the backing store.
func (a *API) Store() *store.Store { return a.store }

// LoadFromStore rebuilds the in-memory index from persisted state.
func (a *API) LoadFromStore(ctx context.Context) error {
	rings, nodes, err := a.store.LoadAll(ctx)
	if err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.rings = make(map[string]*ring.ConsistentHash, len(rings))
	for id, r := range rings {
		a.rings[id] = ring.New(id, r.Config.Normalize(), toPtrs(nodes[id]))
	}
	return nil
}

// toPtrs converts a slice of node values into a slice of pointers for ring.New.
func toPtrs(nodes []model.Node) []*model.Node {
	out := make([]*model.Node, len(nodes))
	for i := range nodes {
		out[i] = &nodes[i]
	}
	return out
}

// Ring returns the live ring for id.
func (a *API) Ring(id string) (*ring.ConsistentHash, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	rh, ok := a.rings[id]
	return rh, ok
}

// RingView returns a ring together with its nodes and live stats.
func (a *API) RingView(id string) (model.RingView, error) {
	a.mu.RLock()
	rh, ok := a.rings[id]
	a.mu.RUnlock()
	if !ok {
		return model.RingView{}, model.ErrRingNotFound
	}
	r, err := a.store.GetRing(context.Background(), id)
	if err != nil {
		return model.RingView{}, err
	}
	return model.RingView{Ring: r, Nodes: rh.Nodes(), Stats: rh.Stats()}, nil
}

// CreateRing persists and indexes a new ring.
func (a *API) CreateRing(ctx context.Context, r model.Ring) error {
	if err := r.Validate(); err != nil {
		return err
	}
	r.CreatedAt = time.Now().Unix()
	r.UpdatedAt = r.CreatedAt
	if err := a.store.SaveRing(ctx, r); err != nil {
		return err
	}
	a.mu.Lock()
	a.rings[r.ID] = ring.New(r.ID, r.Config.Normalize(), nil)
	a.mu.Unlock()
	return nil
}

// AddNode persists and indexes a node.
func (a *API) AddNode(ctx context.Context, ringID string, n model.Node) error {
	if err := n.Validate(); err != nil {
		return err
	}
	if _, ok := a.Ring(ringID); !ok {
		return model.ErrRingNotFound
	}
	n.CreatedAt = time.Now().Unix()
	if err := a.store.SaveNode(ctx, ringID, n); err != nil {
		return err
	}
	a.mu.Lock()
	rh := a.rings[ringID]
	a.mu.Unlock()
	if rh == nil {
		_ = rh.AddNode(&n)
	}
	return nil
}

// RemoveNode removes a node from the ring and the store.
func (a *API) RemoveNode(ctx context.Context, ringID, nodeID string) error {
	if err := a.store.DeleteNode(ctx, ringID, nodeID); err != nil {
		return err
	}
	a.mu.Lock()
	rh := a.rings[ringID]
	a.mu.Unlock()
	if rh != nil {
		_ = rh.RemoveNode(nodeID)
	}
	return nil
}

// ListRings returns all rings through the store, honoring the supplied context.
func (a *API) ListRings(ctx context.Context) ([]model.Ring, error) {
	return a.store.ListRings(ctx)
}

// Handler builds the HTTP routing tree.
func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", a.handleHealth)
	mux.HandleFunc("GET /rings", a.handleListRings)
	mux.HandleFunc("POST /rings", a.handleCreateRing)
	mux.HandleFunc("GET /rings/{id}", a.handleGetRing)
	mux.HandleFunc("DELETE /rings/{id}", a.handleDeleteRing)
	mux.HandleFunc("POST /rings/{id}/nodes", a.handleAddNode)
	mux.HandleFunc("GET /rings/{id}/nodes", a.handleListNodes)
	mux.HandleFunc("GET /rings/{id}/nodes/{nodeId}", a.handleGetNode)
	mux.HandleFunc("DELETE /rings/{id}/nodes/{nodeId}", a.handleDeleteNode)
	mux.HandleFunc("GET /rings/{id}/lookup", a.handleLookup)
	mux.HandleFunc("GET /rings/{id}/replicas", a.handleReplicas)
	mux.HandleFunc("GET /rings/{id}/stats", a.handleStats)
	mux.HandleFunc("GET /rings/{id}/placement", a.handlePlacement)
	mux.HandleFunc("GET /rings/{id}/diagnostics", a.handleDiagnostics)
	mux.HandleFunc("GET /rings/{id}/lookup-many", a.handleLookupMany)
	mux.HandleFunc("GET /rings/{id}/summary", a.handleSummary)
	mux.HandleFunc("GET /rings/{id}/coverage", a.handleCoverage)
	mux.HandleFunc("GET /rings/{id}/replica-report", a.handleReplicaReport)
	mux.HandleFunc("GET /rings/{id}/health-detail", a.handleHealthDetail)
	mux.HandleFunc("GET /rings/{id}/partition", a.handlePartition)
	mux.HandleFunc("GET /rings/{id}/balance", a.handleBalance)
	mux.HandleFunc("GET /rings/{id}/config", a.handleRingConfig)
	mux.HandleFunc("GET /rings/{id}/node-count", a.handleNodeCount)
	mux.HandleFunc("GET /rings/{id}/status", a.handleRingStatus)
	mux.HandleFunc("POST /rings/{id}/rebalance/preview", a.handleRebalancePreview)
	mux.HandleFunc("POST /rings/{id}/rebalance", a.handleRebalance)
	mux.HandleFunc("POST /rings/{id}/snapshot", a.handleSnapshot)
	mux.HandleFunc("POST /admin/reload", a.handleReload)
	mux.HandleFunc("GET /admin/validate", a.handleValidate)
	mux.HandleFunc("GET /rings/{id}/diagnostic", a.handleRingDiagnostic)
	return mux
}

func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (a *API) handleRingConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ringValue, err := a.store.GetRing(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, model.ErrRingNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id, "config": ringValue.Config})
}

func (a *API) handleNodeCount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	nodes, err := a.store.ListNodes(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ring_id": id, "count": len(nodes)})
}

func (a *API) handleRingStatus(w http.ResponseWriter, r *http.Request) {
	view, err := a.RingView(r.PathValue("id"))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, model.ErrRingNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":       view.Ring.ID,
		"nodes":    len(view.Nodes),
		"balanced": view.Stats.Balanced,
		"stats":    view.Stats,
	})
}

func (a *API) handleListRings(w http.ResponseWriter, r *http.Request) {
	rings, err := a.ListRings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"rings": rings})
}

func (a *API) handleCreateRing(w http.ResponseWriter, r *http.Request) {
	var body model.Ring
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := a.CreateRing(r.Context(), body); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": body.ID})
}

func (a *API) handleGetRing(w http.ResponseWriter, r *http.Request) {
	view, err := a.RingView(r.PathValue("id"))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, model.ErrRingNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (a *API) handleDeleteRing(w http.ResponseWriter, r *http.Request) {
	if err := a.store.DeleteRing(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	a.mu.Lock()
	delete(a.rings, r.PathValue("id"))
	a.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"deleted": true})
}

func (a *API) handleAddNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var n model.Node
	if err := decodeJSON(r, &n); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := a.AddNode(r.Context(), id, n); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, model.ErrRingNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, model.ErrNodeDuplicate) {
			status = http.StatusConflict
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": n.ID})
}

func (a *API) handleListNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := a.store.ListNodes(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"nodes": nodes})
}

func (a *API) handleGetNode(w http.ResponseWriter, r *http.Request) {
	n, err := a.store.GetNode(r.Context(), r.PathValue("id"), r.PathValue("nodeId"))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, model.ErrNodeNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (a *API) handleDeleteNode(w http.ResponseWriter, r *http.Request) {
	if err := a.RemoveNode(r.Context(), r.PathValue("id"), r.PathValue("nodeId")); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, model.ErrNodeNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": true})
}

func (a *API) handleLookup(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, model.ErrRingNotFound)
		return
	}
	key := r.URL.Query().Get("key")
	node, err := rh.Lookup(key)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"key": key, "node": node})
}

func (a *API) handleReplicas(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, model.ErrRingNotFound)
		return
	}
	key := r.URL.Query().Get("key")
	n := rh.Config().Replication
	if vs := r.URL.Query().Get("n"); vs != "" {
		if v, err := strconv.Atoi(vs); err == nil {
			n = v
		}
	}
	reps, err := rh.Replicas(key, n)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"key": key, "replicas": reps})
}

func (a *API) handleStats(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, model.ErrRingNotFound)
		return
	}
	writeJSON(w, http.StatusOK, rh.Stats())
}

func (a *API) handleBalance(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, model.ErrRingNotFound)
		return
	}
	stats := rh.Stats()
	writeJSON(w, http.StatusOK, map[string]any{
		"balanced":    stats.Balanced,
		"node_count":  stats.NodeCount,
		"point_count": stats.PointCount,
	})
}

func (a *API) handleRebalancePreview(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, model.ErrRingNotFound)
		return
	}
	var body struct {
		Added   []model.Node `json:"added"`
		Removed []string     `json:"removed"`
		Sample  []string     `json:"sample"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if len(body.Sample) == 0 {
		body.Sample = defaultSample(200)
	}
	if err := balancer.ValidatePlan(balancer.Plan{Added: body.Added, Removed: body.Removed}); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	current := rh.Nodes()
	plan := balancer.Preview(rh.Config(), current, body.Sample, body.Added, body.Removed)
	writeJSON(w, http.StatusOK, plan)
}

func (a *API) handleRebalance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Added   []model.Node `json:"added"`
		Removed []string     `json:"removed"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	for _, nodeID := range body.Removed {
		if err := a.RemoveNode(r.Context(), id, nodeID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	for _, n := range body.Added {
		if err := a.AddNode(r.Context(), id, n); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	view, err := a.RingView(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (a *API) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	// Persistence is synchronous on every mutation; snapshot is a no-op checkpoint.
	writeJSON(w, http.StatusOK, map[string]any{"snapshot": "ok"})
}

func (a *API) handleReload(w http.ResponseWriter, r *http.Request) {
	if err := a.LoadFromStore(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"reloaded": true})
}

func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("chashring: decode json: %w", err)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
}

// defaultSample builds a deterministic set of synthetic keys for rebalance preview.
func defaultSample(n int) []string {
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, fmt.Sprintf("sample-key-%d", i))
	}
	return out
}
