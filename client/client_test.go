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

func TestNewClient_returnsErrorWhenGivenUnparseableURL(t *testing.T) {
	_, err := NewClient("this is not a host", nil)
	assert.Equal(t, "input error - invalid host", err.Error())
}

func TestNewClient_noScheme(t *testing.T) {
	_, err := NewClient("0.0.0.0:8080", nil)
	assert.Equal(t, `input error - invalid host: parse "0.0.0.0:8080": first path segment in URL cannot contain colon`, err.Error())
}

func TestNewClient_removesPathFromHost(t *testing.T) {
	c, err := NewClient("http://0.0.0.0:8080/too/much/path", nil)
	assert.NoError(t, err)
	assert.Equal(t, "http://0.0.0.0:8080", c.host)
}

func TestNewClient_returnsCorrectlyConfiguredClient(t *testing.T) {
	// check that client configures correctly with injected transport
	c, err := NewClient("http://0.0.0.0:8080", &mockRoundTripper{})
	assert.NoError(t, err)
	assert.Equal(t, "http://0.0.0.0:8080", c.host)

	expectedTransport := &requiredHeadersTransportDecorator{
		host: c.host, transport: &mockRoundTripper{},
	}
	assert.Equal(t, expectedTransport, c.httpClient.Transport)

	// check that client defaults to http.DefaultTransport when no custom transport is injected
	c, err = NewClient("http://0.0.0.0:8080", nil)
	assert.NoError(t, err)
	assert.Equal(t, "http://0.0.0.0:8080", c.host)

	expectedTransport = &requiredHeadersTransportDecorator{host: c.host, transport: http.DefaultTransport}
	assert.Equal(t, expectedTransport, c.httpClient.Transport)
}

func TestClient_AddsRequiredHeaders(t *testing.T) {
	// setup mockRoundTripper that checks outgoing request for required headers
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "application/vnd.api+json", req.Header.Get("Accept"))
			assert.Equal(t, "http://0.0.0.0:8080", req.Header.Get("Host"))

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

	c, err := NewClient("http://0.0.0.0:8080", mrt)
	assert.NoError(t, err)

	// check GET request has expected headers
	req, err := http.NewRequest(http.MethodGet, "http://0.0.0.0:8080/this/is/a/fake", nil)
	assert.NoError(t, err)

	_, err = c.httpClient.Do(req)
	assert.NoError(t, err)

	// check POST request has expected headers including those only required for requests with body
	req, err = http.NewRequest(
		http.MethodPost,
		"http://0.0.0.0:8080/this/is/a/fake",
		ioutil.NopCloser(bytes.NewBufferString("this request has a body")),
	)
	assert.NoError(t, err)

	_, err = c.httpClient.Do(req)
	assert.NoError(t, err)

}

// the following testing could also be done with the mocked roundtripper but I thought I would throw this in
// for purely demo purposes as a nice way to test request issuing code

func TestClientget_checkRequestArrivesAndReturnsExpectedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/this/is/a/fake", r.URL.String())

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"this": "is the response"}`))
		assert.NoError(t, err)
	}))

	defer server.Close()

	c, err := NewClient(server.URL, nil)
	assert.NoError(t, err)

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
		_, err = w.Write([]byte(`{"this": "is the response"}`))
		assert.NoError(t, err)
	}))

	defer server.Close()

	c, err := NewClient(server.URL, nil)
	assert.NoError(t, err)

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

	c, err := NewClient(server.URL, nil)
	assert.NoError(t, err)

	resp, err := c.delete(context.Background(), "/this/is/a/fake")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
