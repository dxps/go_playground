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
	"reflect"
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

	// TODO: Extend testing with multiple cases, based on the same Table Driven Tests approach.

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

func Test_getSecretHandler(t *testing.T) {

	// Setup
	f, a, err := setup()
	t.Cleanup(func() { cleanupTest(f) })
	if err != nil {
		t.Fatalf("Setup of the test failed: %v", err)
	}

	// Adding a secret first, before getting it.
	inp := AddSecretInput{PlainText: "snap1"}
	rBody, _ := json.Marshal(inp)
	body := bytes.NewBuffer(rBody)
	r := httptest.NewRequest(http.MethodPost, "http://localhost/secrets", body)
	w := httptest.NewRecorder()
	a.addSecretHandler(w, r)

	// Execute

	testcases := []struct {
		description     string
		requestMethod   string
		requestURLPath  string
		expResponseCode int
		expResponseBody any
	}{
		{
			description:     "When wrong (not GET) HTTP method",
			requestMethod:   http.MethodPost,
			requestURLPath:  "http://localhost/secrets/{not-relevant}",
			expResponseCode: http.StatusBadRequest,
			expResponseBody: nil,
		},
		{
			description:     "When no ID provided",
			requestMethod:   http.MethodGet,
			requestURLPath:  "http://localhost/secrets/",
			expResponseCode: http.StatusBadRequest,
			expResponseBody: ResponseError{Error: "ID (part of '/secrets/{ID}' URL path) is missing."},
		},
		{
			description:     "When correct (previously added) ID provided",
			requestMethod:   http.MethodGet,
			requestURLPath:  "http://localhost/secrets/57c3c3e8874431efd7ae79a8972bdfd4", // ID corresponding to "snap1" secret.
			expResponseCode: http.StatusOK,
			expResponseBody: NewGetSecretOutput("snap1"),
		},
	}

	for _, tc := range testcases {
		r := httptest.NewRequest(tc.requestMethod, tc.requestURLPath, nil)
		w := httptest.NewRecorder()

		// Execute
		a.getSecretHandler(w, r)

		// Evaluate
		if tc.expResponseCode != w.Result().StatusCode {
			t.Fatalf(fmt.Sprintf("%s: on status code expected %d, got %d", tc.description, tc.expResponseCode, w.Result().StatusCode))
		}

		if tc.expResponseBody == nil {
			// It's not part of the test, so let's continue with the next test.
			continue
		}
		var got any
		var err error
		switch rbt := tc.expResponseBody.(type) {
		case ResponseError:
			got, err = NewResponseErrorFromBytes(w.Body.Bytes())
		case GetSecretOutput:
			got, err = NewGetSecretOutputFromBytes(w.Body.Bytes())
		default:
			t.Fatalf("Got unexpected '%v' as response body type", rbt)
		}

		if err != nil {
			t.Fatalf(fmt.Sprintf("%s: failed to unmarshal response body: %v", tc.description, err))
		}
		if reflect.DeepEqual(tc.expResponseBody, got) {
			t.Fatalf(fmt.Sprintf("%s: on body expected %v, got %v", tc.description, tc.expResponseBody, got))
		}
	}
}
