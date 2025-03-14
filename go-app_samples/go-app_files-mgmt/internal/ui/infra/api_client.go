//go:build js

package infra

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"go-app_files-mgmt/internal/ui/infra/fetch"
)

const Timeout = 5 * time.Second

type ApiClient struct {
	endpoint string
}

// NewApiClient creates a new `ApiClient` instance.
func NewApiClient(endpoint string) *ApiClient {
	return &ApiClient{
		endpoint,
	}
}

func (apic *ApiClient) Post(path string, reqBody []byte) (respBody []byte, err error) {
	return apic.doRequest(context.Background(), "POST", path, "", reqBody)
}

func (apic *ApiClient) Get(path string) (respBody []byte, err error) {
	return apic.doRequest(context.Background(), "GET", path, "", nil)
}

func (c *ApiClient) SendFile(path string, contentTypeHeader string, contentBytes []byte) (respBody []byte, err error) {

	return c.doRequest(context.Background(), "POST", path, contentTypeHeader, contentBytes)
}

func (c *ApiClient) doRequest(
	ctx context.Context,
	method, path string,
	contentType string,
	reqBody []byte,
) (respBody []byte, err error) {

	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	// u, err := url.Parse(c.endpoint + path)
	// if err != nil {
	// return nil, fmt.Errorf("failed to parse request URL: %w", err)
	// }
	// Encode the query params.
	// q := u.Query()
	// for k, v := range queryParams {
	// q.Add(k, v)
	// }
	// u.RawQuery = q.Encode()

	// req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(reqBody))
	// if err != nil {
	// return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	// }
	//
	// req.Header.Set("Content-Type", "application/json")
	// Perform the HTTP request.
	// resp, err := c.client.Do(req)

	resp, err := fetch.Fetch(c.endpoint+path, &fetch.Opts{
		Method: method,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
		Body:   bytes.NewBuffer(reqBody),
		Signal: ctx,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to do the HTTP request: %w", err)
	}
	// defer resp.Body.Close()

	// respBody, err = io.ReadAll(resp.Body)
	// if err != nil {
	// return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	// }

	// check if status is OK
	// if resp.StatusCode != http.StatusOK {
	// return nil, fmt.Errorf("server responded with status code: %d: %w",
	// resp.StatusCode, errors.New(string(respBody)))
	// }

	if resp.Status != 200 {
		return nil, fmt.Errorf("server responded with status code: %v: %w",
			resp.Status, errors.New(string(resp.Body)))
	}

	return resp.Body, nil
}
