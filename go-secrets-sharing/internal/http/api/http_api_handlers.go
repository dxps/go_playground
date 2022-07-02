package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/errors"
	log "github.com/sirupsen/logrus"
)

func (a *HttpApi) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}
	// Blindly respond OK, for now.
	// Ignore the HTTP method used or the current app state.
	w.WriteHeader(http.StatusOK)
}

func (a *HttpApi) addSecretHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
	}
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Request body read error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	var input AddSecretInput
	if err := json.Unmarshal(rBody, &input); err != nil {
		log.Errorf("Json unmarshall of input error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	hash := a.secrets.Store(input.PlainText)
	respBody := NewAddSecretOutput(hash)
	respData, err := json.Marshal(respBody)
	if err != nil {
		log.Errorf("[secretsHandler] json.Marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(respData)
}

func (a *HttpApi) getSecretHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	}
	hash := strings.TrimPrefix(r.URL.Path, "/secrets/")
	if len(strings.TrimSpace(hash)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	secret, err := a.secrets.Retrieve(hash)
	if err != nil {
		if err == errors.EntryNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Errorf("[getSecretHandler] Got error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	respBody := NewGetSecretOutput(secret)
	respData, err := json.Marshal(respBody)
	if err != nil {
		log.Errorf("[getSecretHandler] json.Marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(respData)
}
