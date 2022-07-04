package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/OJOMB/form3-fake-account-client/accounts"
	"github.com/stretchr/testify/assert"
)

// TestFetch_return200WithValidRespBody_SuccessPath tests the happy path of the Fetch function by simulating
// a successful response from the API
// All fields including deprecated fields are returned to the client and checked for in the client response
func TestFetch_return200WithValidRespBody_SuccessPath(t *testing.T) {
	// mock roundtripper will return a response with all fields including deprecated
	// this is of course unfaithful to the actual API because of server-side validations that would result in this never happening
	// despite this doing it this way still serves to allow us to test the clients ability to accurately gather all fields in the response
	respBody := fmt.Sprintf(
		`{
			"data": %s,
			"links": {
				"self": "/v1/organisation/accounts/1dfaf917-c6d6-4e18-b7e7-972e66492976"
			}
		}`,
		getTestDataAccountDataAllFields(getDummyTime().Format(time.RFC3339), getDummyTime().Format(time.RFC3339), 0),
	)

	// roundtripper returns 200 success response with above body
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(respBody)),
			}, nil
		},
	}

	client := NewClient("http://localhost:8080", mrt)
	resp, err := client.Fetch(context.Background(), "1dfaf917-c6d6-4e18-b7e7-972e66492976")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	dummyTime := ptrTime(getDummyTime())
	expectedAccountData := accounts.AccountData{
		Type:           "accounts",
		ID:             "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		OrganisationID: "caca9817-6936-4da4-96e7-9ce93206070f",
		CreatedOn:      dummyTime,
		ModifiedOn:     dummyTime,
		Version:        ptrInt64(0),
		Attributes: &accounts.AccountAttributes{
			AccountClassification:   accounts.AccountClassificationPersonal,
			AccountNumber:           "10000004",
			BankID:                  "400302",
			BankIDCode:              "GBDSC",
			BaseCurrency:            "GBP",
			Bic:                     "NWBKGB42",
			Country:                 ptrStr("GB"),
			CustomerID:              "12345",
			Iban:                    "GB28NWBK40030212764204",
			Name:                    []string{"Jane Doe"},
			NameMatchingStatus:      accounts.AccountNameMatchingStatusOptedOut,
			AlternativeNames:        []string{"Sam Holder"},
			JointAccount:            ptrBool(false),
			SecondaryIdentification: "A1B2C3D4",
			Status:                  ptrAccountStatus(accounts.AccountStatusConfirmed),
			StatusReason:            "unspecified",
			UserDefinedData: []accounts.UserDefinedData{
				{
					Key:   "Some account related key",
					Value: "Some account related value",
				},
			},
			ValidationType:              "card",
			ReferenceMask:               "############",
			AcceptanceQualifier:         "same_day",
			Title:                       "Mrs",
			FirstName:                   "Jane",
			BankAccountName:             "Jane Doe",
			AlternativeBankAccountNames: []string{"Sam Holder"},
			ProcessingService:           "processing_service_1",
			UserDefinedInformation:      "Some account related value",
			AccountMatchingOptOut:       false,
			Switched:                    ptrBool(false),
		},
	}

	expectedresp := &accounts.Response{
		Data: &expectedAccountData,
		Links: &accounts.Links{
			Self: "/v1/organisation/accounts/1dfaf917-c6d6-4e18-b7e7-972e66492976",
		},
	}

	assert.EqualValues(t, expectedresp, resp)
}

func TestFetch_return200WithInvalidJsonRespBody_FailurePath(t *testing.T) {
	respBody := `{this is invalid JSON": {}}`

	// roundtripper returns 200 success response with above body
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       ioutil.NopCloser(bytes.NewBufferString(respBody)),
			}, nil
		},
	}

	client := NewClient("http://localhost:8080", mrt)

	resp, err := client.Fetch(context.Background(), "1dfaf917-c6d6-4e18-b7e7-972e66492976")
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, `internal error - failed to unmarshal response body: invalid character 't' looking for beginning of object key string`, err.Error())
}

func TestFetch_returnUnreadableBody_FailurePath(t *testing.T) {
	// roundtripper returns unreadable response
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusCreated,
				Body:       ioutil.NopCloser(errReader(0)),
			}

			return resp, nil
		},
	}

	client := NewClient("http://localhost:8080", mrt)
	resp, err := client.Fetch(context.Background(), "1dfaf917-c6d6-4e18-b7e7-972e66492976")
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, "internal error - failed to read response body: failed to read", err.Error())
}

func TestFetch_FailedSend_FailurePath(t *testing.T) {
	// roundtripper returns erroneous response
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("nope")
		},
	}

	client := NewClient("http://localhost:8080", mrt)
	resp, err := client.Fetch(context.Background(), "1dfaf917-c6d6-4e18-b7e7-972e66492976")
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, `internal error - failed to send http request: Get "http://localhost:8080/v1/organisation/accounts/1dfaf917-c6d6-4e18-b7e7-972e66492976": nope`, err.Error())
}

func TestCreate_returnNon200_FailurePath(t *testing.T) {
	testCases := []struct {
		name         string
		serverErrMsg string
		statusCode   int
	}{
		{
			name:         "server returns 400 Bad Request",
			serverErrMsg: "bad request",
			statusCode:   http.StatusBadRequest,
		},
		{
			name:         "server returns 500 Internal Server Error",
			serverErrMsg: "internal server error",
			statusCode:   http.StatusInternalServerError,
		},
		{
			name:         "server returns 404 Not Found",
			serverErrMsg: "not found",
			statusCode:   http.StatusNotFound,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%s", idx, tc.name), func(t *testing.T) {
			respBody := fmt.Sprintf(`{"error_message": "%s"}`, tc.serverErrMsg)
			mrt := &mockRoundTripper{
				transportFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tc.statusCode,
						Body:       ioutil.NopCloser(bytes.NewBufferString(respBody)),
					}, nil
				},
			}
			client := NewClient("http://localhost:8080", mrt)
			resp, err := client.Fetch(context.Background(), "1dfaf917-c6d6-4e18-b7e7-972e66492976")
			assert.Error(t, err)
			assert.Nil(t, resp)

			expectedErrMsg := fmt.Sprintf("api error - failed to fetch account, status code %d: %s", tc.statusCode, tc.serverErrMsg)
			assert.Equal(t, expectedErrMsg, err.Error())
		})
	}
}

func TestFetch_returnNon200ResponseWithInvalidJson_FailurePath(t *testing.T) {
	respBody := `{"error_message": "this is invalid json}`

	// roundtripper returns 500 error response with above body
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBufferString(respBody)),
			}, nil
		},
	}
	client := NewClient("http://localhost:8080", mrt)
	resp, err := client.Fetch(context.Background(), "1dfaf917-c6d6-4e18-b7e7-972e66492976")
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, "internal error - failed to unmarshal response body: unexpected end of JSON input", err.Error())
}

func TestFetch_return200ResponseWithInvalidJson_FailurePath(t *testing.T) {
	respBody := `{this is invalid JSON": {}}`

	// roundtripper returns 500 error response with above body
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(respBody)),
			}, nil
		},
	}
	client := NewClient("http://localhost:8080", mrt)
	resp, err := client.Fetch(context.Background(), "1dfaf917-c6d6-4e18-b7e7-972e66492976")
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, "internal error - failed to unmarshal response body: invalid character 't' looking for beginning of object key string", err.Error())
}

func TestFetch_AccountWithEmptyID_FailurePath(t *testing.T) {
	client := NewClient("http://localhost:8080", nil)
	resp, err := client.Fetch(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, "input error - accountID cannot be empty", err.Error())
}
