package api

import (
	"net/http"
	"task116-chashring/internal/model"
)

func (a *API) handleValidate(w http.ResponseWriter, r *http.Request) {
	if err := a.store.ValidateAll(r.Context()); err != nil {
		writeError(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"valid": true})
}
func (a *API) handleRingDiagnostic(w http.ResponseWriter, r *http.Request) {
	rh, ok := a.Ring(r.PathValue("id"))
	if !ok {
		writeError(w, 404, model.ErrRingNotFound)
		return
	}
	writeJSON(w, 200, rh.Diagnostic())
}
