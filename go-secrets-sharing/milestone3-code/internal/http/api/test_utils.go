package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/domain"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/repo"
)

func setup() (*os.File, *HttpApi, error) {

	f, err := os.CreateTemp("", "TestHandlerHealthcheck_secrets_")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create temp file: %v", err)
	}
	r, err := repo.NewRepo(f.Name(), "testPass", "testSalt")
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
	// log.Printf("[cleanupTest] Removing '%s' file", f.Name())
	_ = f.Close()
	_ = os.Remove(f.Name())
}

type TestCase struct {
	description        string // Description of the test, like "When something is provided".
	requestMethod      string // The `http.Method...` to be used.
	requestURLPath     string // The endpoint to be called.
	requestBody        any    // Any object that gets JSON marshalled.
	expectResponseCode int    // Expected `http.Status...` code
	expectResponseBody any    // Expected body as object that gets JSON unmarshalled from response body bytes.
}

func runHandlerTests(cases []TestCase, handler http.HandlerFunc) error {

	for _, tc := range cases {
		var body io.Reader
		if tc.requestBody != nil {
			bs, err := json.Marshal(tc.requestBody)
			if err != nil {
				return fmt.Errorf("Unable to test since its prep failed: %v", err)
			}
			body = bytes.NewBuffer(bs)
		}
		r := httptest.NewRequest(tc.requestMethod, tc.requestURLPath, body)
		w := httptest.NewRecorder()

		// Execute
		handler(w, r)

		// Evaluate
		var gotBody any
		var testErr error
		if err := json.Unmarshal(w.Body.Bytes(), &gotBody); err != nil {
			testErr = err
		}
		if tc.expectResponseCode != w.Result().StatusCode {
			return fmt.Errorf("%s, expected status code %d, got %d. Other details: gotBody: %v testErr: %v",
				tc.description, tc.expectResponseCode, w.Result().StatusCode, gotBody, testErr)
		}

		if tc.expectResponseBody == nil {
			// No need to evaluate the response, continue with the next test.
			continue
		}

		gotBs := w.Body.Bytes()
		var got any
		var err error
		switch rbt := tc.expectResponseBody.(type) {
		case *ResponseError:
			got, err = NewResponseErrorFromBytes(gotBs)
		case *AddSecretOutput:
			got, err = NewAddSecretOutputFromBytes(gotBs)
		case *GetSecretOutput:
			got, err = NewGetSecretOutputFromBytes(gotBs)
		default:
			return fmt.Errorf("%s, got unexpected '%v' as response body type. Details: gotBody: %v",
				tc.description, rbt, gotBody)
		}

		if err != nil {
			return fmt.Errorf("%s, failed to unmarshal response body: %v", tc.description, err)
		}
		if !reflect.DeepEqual(tc.expectResponseBody, got) {
			return fmt.Errorf("%s, expected body %v, got %v", tc.description, tc.expectResponseBody, got)
		}
	}
	return nil
}
