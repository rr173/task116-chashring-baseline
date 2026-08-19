package api

import (
	"net/http"
	"task116-chashring/internal/model"
)

func validateNode(w http.ResponseWriter, n model.Node) bool {
	if err := n.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return false
	}
	return true
}
