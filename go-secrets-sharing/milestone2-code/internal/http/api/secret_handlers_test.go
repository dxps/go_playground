package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_addSecretHandler(t *testing.T) {

	// Setup
	f, a, err := setup()
	t.Cleanup(func() { cleanupTest(f) })
	if err != nil {
		t.Fatalf("Setup of the test failed: %v", err)
	}

	testcases := []TestCase{
		{
			description:        "When a standard request is provided",
			requestMethod:      http.MethodPost,
			requestURLPath:     "http://localhost/secrets",
			requestBody:        AddSecretInput{PlainText: "snap1"},
			expectResponseCode: http.StatusOK,
			expectResponseBody: NewAddSecretOutput("57c3c3e8874431efd7ae79a8972bdfd4"),
		},
		{
			description:        "When no body is provided",
			requestMethod:      http.MethodPost,
			requestURLPath:     "http://localhost/secrets",
			requestBody:        nil,
			expectResponseCode: http.StatusBadRequest,
			expectResponseBody: NewResponseError("Cannot decode request body: unexpected end of JSON input"),
		},
		{
			description:        "When wrong (PUT, not POST) method is provided",
			requestMethod:      http.MethodPut,
			requestURLPath:     "http://localhost/secrets",
			requestBody:        nil,
			expectResponseCode: http.StatusBadRequest,
			expectResponseBody: NewResponseError("Only POST method can be used."),
		},
	}

	// Execute & Evaluate
	if err := runHandlerTests(testcases, a.addSecretHandler); err != nil {
		t.Fatalf(err.Error())
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
	rBody, err := json.Marshal(inp)
	if err != nil {
		t.Fatalf("Unable to test since its prep failed: %v", err)
	}
	body := bytes.NewBuffer(rBody)
	r := httptest.NewRequest(http.MethodPost, "http://localhost/secrets", body)
	w := httptest.NewRecorder()
	a.addSecretHandler(w, r)

	testcases := []TestCase{
		{
			description:        "When wrong (not GET) HTTP method",
			requestMethod:      http.MethodPost,
			requestURLPath:     "http://localhost/secrets/{not-relevant}",
			expectResponseCode: http.StatusBadRequest,
			expectResponseBody: nil,
		},
		{
			description:        "When no ID provided",
			requestMethod:      http.MethodGet,
			requestURLPath:     "http://localhost/secrets/",
			expectResponseCode: http.StatusBadRequest,
			expectResponseBody: NewResponseError("ID (part of '/secrets/{ID}' URL path) is missing."),
		},
		{
			description:        "When correct (previously added) ID is provided",
			requestMethod:      http.MethodGet,
			requestURLPath:     "http://localhost/secrets/57c3c3e8874431efd7ae79a8972bdfd4", // ID corresponding to "snap1" secret.
			expectResponseCode: http.StatusOK,
			expectResponseBody: NewGetSecretOutput("snap1"),
		},
		{
			description:        "When correct (but already retrieved) ID is provided",
			requestMethod:      http.MethodGet,
			requestURLPath:     "http://localhost/secrets/57c3c3e8874431efd7ae79a8972bdfd4", // ID corresponding to "snap1" secret.
			expectResponseCode: http.StatusNotFound,
			expectResponseBody: nil,
		},
	}

	// Execute & Evaluate
	if err := runHandlerTests(testcases, a.getSecretHandler); err != nil {
		t.Fatalf(err.Error())
	}
}
