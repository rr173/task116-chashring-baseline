package api

import (
	"net/http"
	"task116-chashring/internal/balancer"
	"task116-chashring/internal/model"
)

func (a *API) handleCoverage(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, 404, model.ErrRingNotFound)
		return
	}
	keys := defaultSample(200)
	assignments := make([]string, 0, len(keys))
	for _, key := range keys {
		node, err := rh.Lookup(key)
		if err == nil {
			assignments = append(assignments, node)
		}
	}
	writeJSON(w, 200, balancer.CoverageOf(assignments))
}
