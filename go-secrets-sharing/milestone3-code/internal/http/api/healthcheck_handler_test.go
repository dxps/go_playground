package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
