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

// TestCreate_return201WithValidRespBody_SuccessPath tests the happy path of the Create function by simulating
// a successful response from the API
// All fields including deprecated fields are returned to the client and checked for in the client response
func TestCreate_return201WithValidRespBody_SuccessPath(t *testing.T) {
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

	// roundtripper returns 201 success response with above body
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       ioutil.NopCloser(bytes.NewBufferString(respBody)),
			}, nil
		},
	}

	c, err := NewClient("http://localhost:8080", mrt)
	assert.NoError(t, err)

	account := accounts.AccountData{
		Type:           "accounts",
		ID:             "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		OrganisationID: "caca9817-6936-4da4-96e7-9ce93206070f",
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

	resp, err := c.Create(context.Background(), account)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	expectedresp := &accounts.Response{
		Data: &account,
		Links: &accounts.Links{
			Self: fmt.Sprintf("%s/1dfaf917-c6d6-4e18-b7e7-972e66492976", basev1AccountsPath),
		},
	}
	expectedresp.Data.Version = ptrInt64(0)
	expectedresp.Data.CreatedOn = ptrTime(getDummyTime())
	expectedresp.Data.ModifiedOn = ptrTime(getDummyTime())

	assert.EqualValues(t, expectedresp, resp)
}

func TestCreate_return201WithInvalidJsonRespBody_FailurePath(t *testing.T) {
	respBody := fmt.Sprintf(
		`{
			"data": %s,
			"links": {
				"self": "this is invalid json
			}
		}`,
		getTestDataAccountDataAllFields(getDummyTime().Format(time.RFC3339), getDummyTime().Format(time.RFC3339), 0),
	)

	// roundtripper returns 201 success response with above body
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       ioutil.NopCloser(bytes.NewBufferString(respBody)),
			}, nil
		},
	}

	c, err := NewClient("http://localhost:8080", mrt)
	assert.NoError(t, err)

	account := accounts.AccountData{
		Type:           "accounts",
		ID:             "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		OrganisationID: "caca9817-6936-4da4-96e7-9ce93206070f",
		Attributes: &accounts.AccountAttributes{
			Country: ptrStr("GB"),
			Name:    []string{"Jane Doe"},
		},
	}

	resp, err := c.Create(context.Background(), account)
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, `internal error - failed to unmarshal response body: invalid character '\n' in string literal`, err.Error())
}

func TestCreate_returnUnreadableBody_FailurePath(t *testing.T) {
	// roundtripper returns erroneous response
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusCreated,
				Body:       ioutil.NopCloser(errReader(0)),
			}

			return resp, nil
		},
	}

	c, err := NewClient("http://localhost:8080", mrt)
	assert.NoError(t, err)

	account := accounts.AccountData{
		Type:           "accounts",
		ID:             "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		OrganisationID: "caca9817-6936-4da4-96e7-9ce93206070f",
		Attributes: &accounts.AccountAttributes{
			Country: ptrStr("GB"),
			Name:    []string{"Jane Doe"},
		},
	}

	resp, err := c.Create(context.Background(), account)
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, "internal error - failed to read response body: failed to read", err.Error())
}

func TestCreate_failedRequestSend_FailurePath(t *testing.T) {
	// roundtripper returns erroneous response
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("nope")
		},
	}

	c, err := NewClient("http://localhost:8080", mrt)
	assert.NoError(t, err)

	account := accounts.AccountData{
		Type:           "accounts",
		ID:             "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		OrganisationID: "caca9817-6936-4da4-96e7-9ce93206070f",
		Attributes: &accounts.AccountAttributes{
			Country: ptrStr("GB"),
			Name:    []string{"Jane Doe"},
		},
	}

	resp, err := c.Create(context.Background(), account)
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, `internal error - failed to send http request: Post "http://localhost:8080/v1/organisation/accounts": nope`, err.Error())
}

func TestCreate_returnNon201_FailurePath(t *testing.T) {
	respBody := `{"error_message": "internal server error"}`

	// roundtripper returns 500 error response with above body
	mrt := &mockRoundTripper{
		transportFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBufferString(respBody)),
			}, nil
		},
	}

	c, err := NewClient("http://localhost:8080", mrt)
	assert.NoError(t, err)

	account := accounts.AccountData{
		Type:           "accounts",
		ID:             "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		OrganisationID: "caca9817-6936-4da4-96e7-9ce93206070f",
		Attributes: &accounts.AccountAttributes{
			Country: ptrStr("GB"),
			Name:    []string{"Jane Doe"},
		},
	}

	resp, err := c.Create(context.Background(), account)
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, "api error - failed to create account, status code 500: internal server error", err.Error())
}

func TestCreate_returnNon201ResponseWithInvalidJson_FailurePath(t *testing.T) {
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

	c, err := NewClient("http://localhost:8080", mrt)
	assert.NoError(t, err)

	account := accounts.AccountData{
		Type:           "accounts",
		ID:             "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		OrganisationID: "caca9817-6936-4da4-96e7-9ce93206070f",
		Attributes: &accounts.AccountAttributes{
			Country: ptrStr("GB"),
			Name:    []string{"Jane Doe"},
		},
	}

	resp, err := c.Create(context.Background(), account)
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Equal(t, "internal error - failed to unmarshal response body: unexpected end of JSON input", err.Error())
}
