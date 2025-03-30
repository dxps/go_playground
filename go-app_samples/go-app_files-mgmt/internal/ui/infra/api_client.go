//go:build js && wasm

package infra

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
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
	return apic.doRequest(context.Background(), fetch.MethodPost, path, "", reqBody)
}

func (apic *ApiClient) Get(path string) (respBody []byte, err error) {
	return apic.doRequest(context.Background(), fetch.MethodGet, path, "", nil)
}

func (c *ApiClient) SendFile(
	path string,
	contentTypeHeader string,
	contentBytes []byte,
) (respBody []byte, err error) {
	return c.doRequest(context.Background(), fetch.MethodPost, path, contentTypeHeader, contentBytes)
}

func (c *ApiClient) doRequest(
	ctx context.Context,
	method, path string,
	contentTypeHeader string,
	reqBody []byte,
) (respBody []byte, err error) {

	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	var body io.Reader
	if method == fetch.MethodPost || method == fetch.MethodPut {
		body = bytes.NewBuffer(reqBody)
	}
	resp, err := fetch.Fetch(c.endpoint+path, &fetch.Opts{
		Method: method,
		Headers: map[string]string{
			"Content-Type": contentTypeHeader,
		},
		Body:   body,
		Signal: ctx,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to do the HTTP request: %w", err)
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("Server responded with status code: %v: %w",
			resp.Status, errors.New(string(resp.Body)))
	}

	return resp.Body, nil
}

func (c *ApiClient) DoRequest(
	ctx context.Context,
	method, path string,
	headers map[string]string,
	reqBody []byte,
) (respBody []byte, err error) {

	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	resp, err := fetch.Fetch(c.endpoint+path, &fetch.Opts{
		Method:  method,
		Headers: headers,
		Body:    bytes.NewBuffer(reqBody),
		Signal:  ctx,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to do the HTTP request: %w", err)
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("Server responded with status code: %v: %w",
			resp.Status, errors.New(string(resp.Body)))
	}

	return resp.Body, nil
}
