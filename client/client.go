package client

import "net/http"

const basev1AccountsPath = "/v1/organisation/accounts"

type (
	// Client functions to provide a programatic interface to the fake account API via it's methods
	Client struct {
		*http.Client
		baseURL string
	}
)

// NewClient returns a pointer to a new instance of the fake account API client.
// if transport == nil, we use http.DefaultTransport as RoundTripper
func NewClient(baseURL string, transport http.RoundTripper) *Client {
	return &Client{
		Client: &http.Client{
			Transport: transport,
		},
		baseURL: baseURL,
	}
}
