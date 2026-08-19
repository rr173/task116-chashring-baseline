package api

import (
	"net/http"
	"task116-chashring/internal/balancer"
	"task116-chashring/internal/model"
)

func (a *API) handleHealthDetail(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, 404, model.ErrRingNotFound)
		return
	}
	stats := rh.Stats()
	writeJSON(w, 200, map[string]any{"health": rh.Health(), "metrics": rh.Metrics(), "risk": balancer.RiskOf(stats)})
}
