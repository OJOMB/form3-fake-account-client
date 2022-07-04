package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete_return204_SuccessPath(t *testing.T) {
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusNoContent}, nil
		},
	}

	c, err := NewClient("http://0.0.0.0:8080", mrt)
	assert.NoError(t, err)

	err = c.Delete(context.Background(), "1dfaf917-c6d6-4e18-b7e7-972e66492976", 0)
	assert.NoError(t, err)
}

func TestDelete_returnNon204_FailurePath(t *testing.T) {
	testCases := []struct {
		name        string
		statusCode  int
		respBody    string
		expectedErr string
		accountID   string
	}{
		{
			name:        "resource conflict",
			statusCode:  http.StatusConflict,
			respBody:    `{"error_message": "invalid version"}`,
			expectedErr: "api error - failed to delete account, status code 409: invalid version",
			accountID:   "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		},
		{
			name:        "bad request",
			statusCode:  http.StatusBadRequest,
			respBody:    `{"error_message": "id is not a valid uuid"}`,
			expectedErr: "api error - failed to delete account, status code 400: id is not a valid uuid",
			accountID:   "invalid-uuid",
		},
		{
			name:        "resource not found",
			statusCode:  http.StatusNotFound,
			expectedErr: "api error - failed to delete account, status code 404: account not found",
			accountID:   "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		},
		{
			name:        "internal server error",
			statusCode:  http.StatusInternalServerError,
			expectedErr: "api error - failed to delete account, status code 500: server error",
			accountID:   "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %d", idx+1, tc.statusCode), func(t *testing.T) {
			mrt := &mockRoundTripper{
				transportFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tc.statusCode,
						Body:       ioutil.NopCloser(bytes.NewBufferString(tc.respBody)),
					}, nil
				},
			}

			c, err := NewClient("http://0.0.0.0:8080", mrt)
			assert.NoError(t, err)

			err = c.Delete(context.Background(), tc.accountID, 0)
			assert.Equal(t, tc.expectedErr, err.Error())
		})
	}
}

func TestDelete_returnNon204WithBadJSONinResponseBody(t *testing.T) {
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusConflict,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error_message": "invalid version"`)),
			}, nil
		},
	}

	c, err := NewClient("http://0.0.0.0:8080", mrt)
	assert.NoError(t, err)

	err = c.Delete(context.Background(), "1dfaf917-c6d6-4e18-b7e7-972e66492976", 0)
	assert.Equal(t, "internal error - failed to unmarshal response body: unexpected end of JSON input", err.Error())
}

func TestDelete_returnUnreadableBody_FailurePath(t *testing.T) {
	// roundtripper returns unreadable response
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusConflict,
				Body:       ioutil.NopCloser(errReader(0)),
			}

			return resp, nil
		},
	}

	c, err := NewClient("http://0.0.0.0:8080", mrt)
	assert.NoError(t, err)

	err = c.Delete(context.Background(), "1dfaf917-c6d6-4e18-b7e7-972e66492976", 0)
	assert.Equal(t, "internal error - failed to read response body: failed to read", err.Error())
}

func TestDelete_EmptyAccountID_FailurePath(t *testing.T) {
	c, err := NewClient("http://0.0.0.0:8080", nil)
	assert.NoError(t, err)

	err = c.Delete(context.Background(), "", 0)
	assert.Equal(t, "input error - accountID cannot be empty", err.Error())
}
