package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_returnsCorrectlyConfiguredClient(t *testing.T) {
	// check that client configures correctly with injected transport
	c := NewClient("host", &mockRoundTripper{})
	assert.Equal(t, "host", c.host)

	expectedTransport := &requiredHeadersTransportDecorator{
		host: c.host, transport: &mockRoundTripper{},
	}
	assert.Equal(t, expectedTransport, c.httpClient.Transport)

	// check that client defaults to http.DefaultTransport when no custom transport is injected
	c = NewClient("host", nil)
	assert.Equal(t, "host", c.host)

	expectedTransport = &requiredHeadersTransportDecorator{host: c.host, transport: http.DefaultTransport}
	assert.Equal(t, expectedTransport, c.httpClient.Transport)
}

func TestClient_AddsRequiredHeaders(t *testing.T) {
	// setup mockRoundTripper that checks outgoing request for required headers
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "application/vnd.api+json", req.Header.Get("Accept"))
			assert.Equal(t, "example.com", req.Header.Get("Host"))

			headerDateVal := req.Header.Get("Date")
			assert.NotEmpty(t, headerDateVal)
			_, err := time.Parse(time.RFC3339, headerDateVal)
			assert.NoError(t, err)

			if req.Body != nil {
				assert.Equal(t, "application/vnd.api+json", req.Header.Get("Content-Type"))
				assert.NotEmpty(t, "application/vnd.api+json", req.Header.Get("Content-Length"))
			}

			return &http.Response{}, nil
		},
	}

	c := NewClient("example.com", mrt)

	// check GET request has expected headers
	req, err := http.NewRequest(http.MethodGet, "http://localhost:1234/this/is/a/fake", nil)
	assert.NoError(t, err)
	c.httpClient.Do(req)

	// check POST request has expected headers including those only required for requests with body
	req, err = http.NewRequest(
		http.MethodPost,
		"http://localhost:1234/this/is/a/fake",
		ioutil.NopCloser(bytes.NewBufferString("this request has a body")),
	)
	assert.NoError(t, err)
	c.httpClient.Do(req)
}

// the following testing could also be done with the mocked roundtripper but I thought I would throw this in
// for purely demo purposes as a nice way to test request issuing code

func TestClientget_checkRequestArrivesAndReturnsExpectedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/this/is/a/fake", r.URL.String())

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"this": "is the response"}`))
	}))

	defer server.Close()

	c := NewClient(server.URL, nil)

	resp, err := c.get(context.Background(), "/this/is/a/fake")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{"this": "is the response"}`, string(respBody))
}

func TestClientpost_checkRequestArrivesAndReturnsExpectedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/this/is/a/fake", r.URL.String())

		reqBody, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"this": "is the request"}`, string(reqBody))

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"this": "is the response"}`))
	}))

	defer server.Close()

	c := NewClient(server.URL, nil)

	resp, err := c.post(context.Background(), "/this/is/a/fake", []byte(`{"this": "is the request"}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{"this": "is the response"}`, string(respBody))
}

func TestClientdelete_checkRequestArrivesAndReturnsExpectedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/this/is/a/fake", r.URL.String())

		w.WriteHeader(http.StatusNoContent)
	}))

	defer server.Close()

	c := NewClient(server.URL, nil)

	resp, err := c.delete(context.Background(), "/this/is/a/fake")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
