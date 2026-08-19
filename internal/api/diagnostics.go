package api

import (
	"net/http"
	"task116-chashring/internal/balancer"
)

func (a *API) handleDiagnostics(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	view, err := a.RingView(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	health, err := a.store.Diagnostics(r.Context())
	if err != nil {
		writeError(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"ring": view, "health": health, "heatmap": balancer.Heatmap(view.Stats)})
}
