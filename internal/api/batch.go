package api

import (
	"net/http"
	"strings"
	"task116-chashring/internal/model"
)

func (a *API) handleLookupMany(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, 404, model.ErrRingNotFound)
		return
	}
	keys := []string{}
	if value := r.URL.Query().Get("keys"); value != "" {
		keys = strings.Split(value, ",")
	}
	writeJSON(w, 200, map[string]any{"results": rh.LookupMany(keys)})
}
func (a *API) handleSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := a.store.Summary(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, 404, err)
		return
	}
	writeJSON(w, 200, summary)
}
