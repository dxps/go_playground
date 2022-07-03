package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/apperrs"
)

func (a *HttpApi) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		a.respondError(w, "Only GET method can be used.", http.StatusBadRequest)
		return
	}
	// Blindly respond OK, for now.
	w.WriteHeader(http.StatusOK)
}

func (a *HttpApi) addSecretHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		a.respondError(w, "Only POST method can be used.", http.StatusBadRequest)
		return
	}

	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		a.respondError(w, fmt.Sprintf("Cannot decode request body: %v", err))
		return
	}

	var input AddSecretInput
	if err := json.Unmarshal(rBody, &input); err != nil {
		a.respondError(w, fmt.Sprintf("Cannot decode request body: %v", err), http.StatusBadRequest)
		return
	}

	hash, err := a.secrets.Store(input.PlainText)
	if err != nil {
		a.respondError(w, fmt.Sprintf("Store error: %v", err))
		return
	}

	respBody := NewAddSecretOutput(hash)
	a.respondJSON(w, respBody, http.StatusOK)
}

func (a *HttpApi) getSecretHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		a.respondError(w, "Only GET method can be used.", http.StatusBadRequest)
		return
	}

	hash := strings.TrimPrefix(r.URL.Path, "/secrets/")
	if len(strings.TrimSpace(hash)) == 0 {
		a.respondError(w, "ID (the param, part of '/secrets/{ID}' URL path) is missing.", http.StatusBadRequest)
		return
	}
	secret, err := a.secrets.Retrieve(hash)
	if err != nil {
		if err == apperrs.ErrEntryNotFound {
			a.respondError(w, "Unknown ID.", http.StatusNotFound)
		} else {
			a.respondError(w, err.Error())
		}
		return
	}
	respBody := NewGetSecretOutput(secret)
	a.respondJSON(w, respBody, http.StatusOK)
}
