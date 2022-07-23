package api

import (
	"encoding/json"
	"net/http"
)

func (a *HttpApi) respondJSON(w http.ResponseWriter, data any, status int, headers ...http.Header) {

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *HttpApi) respondError(w http.ResponseWriter, err string, status ...int) {

	// Silent for now, maybe later with some logging verbosity control.
	// log.Printf("Error: %v\n", err)

	respStatus := http.StatusInternalServerError
	if len(status) > 0 {
		respStatus = status[0]
	}
	payload := struct {
		Error string `json:"error"`
	}{
		Error: err,
	}
	a.respondJSON(w, payload, respStatus)
}
