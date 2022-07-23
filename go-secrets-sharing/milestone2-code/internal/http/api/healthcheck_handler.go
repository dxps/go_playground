package api

import "net/http"

func (a *HttpApi) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		a.respondError(w, "Only GET method can be used.", http.StatusBadRequest)
		return
	}
	// Blindly respond OK, for now.
	w.WriteHeader(http.StatusOK)
}
