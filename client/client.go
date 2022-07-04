package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

const basev1AccountsPath = "/v1/organisation/accounts"

type (
	// Client functions to provide a programatic interface to the fake account API via it's methods
	Client struct {
		httpClient *http.Client
		host       string
	}
)

// NewClient returns a pointer to a new instance of the fake account API client.
// if transport == nil, we use http.DefaultTransport as RoundTripper
func NewClient(host string, transport http.RoundTripper) *Client {
	if transport == nil {
		transport = http.DefaultTransport
	}

	// decorate transport to add required headers functionality
	transport = &requiredHeadersTransportDecorator{
		host:      host,
		transport: transport,
	}

	return &Client{
		httpClient: &http.Client{
			Transport: transport,
		},
		host: host,
	}
}

// requiredHeadersTransportDecorator is a custom RoundTripper that decorates the RoundTripper in its transport field
// with the functionality to inject the required headers for the fake account API to requests pre-transport
type requiredHeadersTransportDecorator struct {
	transport http.RoundTripper
	host      string
}

// RoundTrip adds headers required by the fake account API and then hands off to the underlying transport
// https://api-docs.form3.tech/api.html#introduction-and-api-conventions-headers
func (drt *requiredHeadersTransportDecorator) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Host", drt.host)
	req.Header.Add("Date", time.Now().Format(time.RFC3339))

	// only add when the request has a body
	if req.Body != nil {
		req.Header.Add("Content-Type", "application/vnd.api+json")

	}

	return drt.transport.RoundTrip(req)
}

// get creates and sends an HTTP GET request
func (c *Client) get(ctx context.Context, path string) (*http.Response, error) {
	return c.createAndDo(ctx, path, http.MethodGet, nil)
}

// post creates and sends an HTTP POST request
func (c *Client) post(ctx context.Context, path string, body []byte) (*http.Response, error) {
	return c.createAndDo(ctx, path, http.MethodPost, body)
}

// delete creates and sends an HTTP DELETE request
func (c *Client) delete(ctx context.Context, path string) (*http.Response, error) {
	return c.createAndDo(ctx, path, http.MethodDelete, nil)
}

// CreateAndDo is a helper function that creates a request and then calls the Do method on the httpClient
func (c *Client) createAndDo(ctx context.Context, path, method string, body []byte) (*http.Response, error) {
	// create request
	url := fmt.Sprintf("%s%s", c.host, path)
	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, newInternalError("failed to create http request", err)
	}

	// send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, newInternalError("failed to send http request", err)
	}

	return resp, nil
}
