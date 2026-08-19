package api

import (
	"net/http"
	"strconv"
	"task116-chashring/internal/model"
)

func (a *API) handlePlacement(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, 404, model.ErrRingNotFound)
		return
	}
	n := rh.Config().Replication
	if value := r.URL.Query().Get("n"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			n = parsed
		}
	}
	placement, err := rh.Placement(r.URL.Query().Get("key"), n)
	if err != nil {
		writeError(w, 500, err)
		return
	}
	writeJSON(w, 200, placement)
}
