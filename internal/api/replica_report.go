package api

import (
	"net/http"
	"task116-chashring/internal/model"
)

func (a *API) handleReplicaReport(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, 404, model.ErrRingNotFound)
		return
	}
	report, err := rh.ReplicaReport(r.URL.Query().Get("key"), rh.Config().Replication)
	if err != nil {
		writeError(w, 500, err)
		return
	}
	writeJSON(w, 200, report)
}
