package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/domain"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/repo"
)

func setup() (*os.File, *HttpApi, error) {

	f, err := os.CreateTemp("", "TestHandlerHealthcheck_secrets_")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create temp file: %v", err)
	}
	r, err := repo.NewRepo(f.Name())
	if err != nil {
		return f, nil, fmt.Errorf("Failed to Setup the repository: %v", err)
	}
	s := domain.NewSecrets(r)
	a := NewHttpApi(s)
	return f, &a, nil
}

func cleanupTest(f *os.File) {

	if f == nil {
		return
	}
	log.Printf("[cleanupTest] Removing '%s' file", f.Name())
	_ = f.Close()
	_ = os.Remove(f.Name())
}

func Test_healthcheckHandler(t *testing.T) {

	// Setup
	f, a, err := setup()
	t.Cleanup(func() { cleanupTest(f) })
	if err != nil {
		t.Fatalf("Setup of the test failed: %v", err)
	}

	r := httptest.NewRequest(http.MethodGet, "http://localhost/healthcheck", nil)
	w := httptest.NewRecorder()

	// Execute
	a.healthcheckHandler(w, r)

	// Evaluate
	if exp, got := http.StatusOK, w.Result().StatusCode; exp != got {
		t.Fatalf("On http status code, expected %d, got %d", exp, got)
	}

}

func Test_addSecretHandler(t *testing.T) {

	// Setup
	f, a, err := setup()
	t.Cleanup(func() { cleanupTest(f) })
	if err != nil {
		t.Fatalf("Setup of the test failed: %v", err)
	}

	inp := AddSecretInput{PlainText: "snap1"}
	exp := AddSecretOutput{ID: "57c3c3e8874431efd7ae79a8972bdfd4"}
	rBody, _ := json.Marshal(inp)
	body := bytes.NewBuffer(rBody)
	r := httptest.NewRequest(http.MethodPost, "http://localhost/secrets", body)
	w := httptest.NewRecorder()

	// Execute
	a.addSecretHandler(w, r)

	// Evaluate
	if exp, got := http.StatusOK, w.Result().StatusCode; exp != got {
		t.Fatalf("On http status code, expected %d, got %d", exp, got)
	}
	respBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Could not process response: %v", err)
	}
	got, err := NewAddSecretOutputFromBytes(respBody)
	if err != nil {
		t.Fatalf("Could not unmarshal response: %v", err)
	}
	if exp != *got {
		t.Fatalf("On response body, expected %v, got %v", exp, *got)
	}
}

func Test_getSecretHandler_no_id_provided_(t *testing.T) {

	// Setup
	f, a, err := setup()
	t.Cleanup(func() { cleanupTest(f) })
	if err != nil {
		t.Fatalf("Setup of the test failed: %v", err)
	}

	r := httptest.NewRequest(http.MethodGet, "http://localhost/secrets/", nil)
	w := httptest.NewRecorder()

	// Execute
	a.getSecretHandler(w, r)

	// Evaluate

	if exp, got := http.StatusBadRequest, w.Result().StatusCode; exp != got {
		t.Fatalf("On status code, expected %d, got %d", exp, got)
	}

	respBody, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Could not process response: %v", err)
	}
	got, err := NewResponseErrorFromBytes(respBody)
	if err != nil {
		t.Fatalf("Could not unmarshal response: %v", err)
	}
	exp := ResponseError{Error: "ID (the param, part of '/secrets/{ID}' URL path) is missing."}
	if exp != *got {
		t.Fatalf("On response body, expected %v, got %v", exp, *got)
	}
}
