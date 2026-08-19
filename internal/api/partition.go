package api

import (
	"net/http"
	"task116-chashring/internal/model"
)

func (a *API) handlePartition(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, 404, model.ErrRingNotFound)
		return
	}
	keys := defaultSample(100)
	writeJSON(w, 200, map[string]any{"partitions": rh.Partition(keys), "keys": model.KeyCount(keys)})
}
