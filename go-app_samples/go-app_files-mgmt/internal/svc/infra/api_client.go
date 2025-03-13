package infra

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const Timeout = 5 * time.Second

type ApiClient struct {
	client   http.Client
	endpoint string
}

// NewApiClient creates a new `ApiClient` instance.
func NewApiClient(endpoint string) *ApiClient {

	client := http.Client{}
	return &ApiClient{
		client,
		endpoint,
	}
}

// Get makes a POST request to the API endpoint, on the provided path.
func (apic *ApiClient) Post(path string, reqBody []byte) (respBody []byte, err error) {
	return apic.doRequest(context.Background(), http.MethodPost, path, nil, reqBody)
}

// Get makes a GET request to the API endpoint, on the provided path.
func (apic *ApiClient) Get(path string) (respBody []byte, err error) {
	return apic.doRequest(context.Background(), http.MethodGet, path, nil, nil)
}

// Get makes a PUT request to the API endpoint, on the provided path.
func (apic *ApiClient) Put(path string, reqBody []byte) (respBody []byte, err error) {
	return apic.doRequest(context.Background(), http.MethodPut, path, nil, reqBody)
}

// Get makes a DELETE request to the API endpoint, on the provided path.
func (apic *ApiClient) Delete(path string) (respBody []byte, err error) {
	return apic.doRequest(context.Background(), http.MethodDelete, path, nil, nil)
}

// doRequest is a generic (and reusable) method for doing an HTTP request.
func (c *ApiClient) doRequest(
	ctx context.Context,
	method, path string,
	queryParams map[string]string,
	reqBody []byte,
) (respBody []byte, err error) {

	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	u, err := url.Parse(c.endpoint + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request URL: %w", err)
	}
	// Encode the query params.
	q := u.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// Perform the HTTP request.
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do the HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	// check if status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status code: %d: %w",
			resp.StatusCode, errors.New(string(respBody)))
	}

	return respBody, nil
}
