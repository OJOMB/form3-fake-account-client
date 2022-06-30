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

// TestCreate_SuccessPath tests the happy path of the Create function by simulating a successful response from the API
func TestCreate_SuccessPath(t *testing.T) {
	// mock roundtripper will return a response with all fields including deprecated
	// this is of course unfaithful to the actual API because of server-side validations that would result in this never happening
	// despite this it serves to allow us to test the clients ability to accurately gather all fields
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

	client := NewClient("http://localhost:8080", mrt)

	account := accounts.AccountData{
		Type:           "accounts",
		ID:             "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		OrganisationID: "caca9817-6936-4da4-96e7-9ce93206070f",
		Attributes: &accounts.AccountAttributes{
			AccountClassification:   ptrStr("Personal"),
			AccountNumber:           "10000004",
			BankID:                  "400302",
			BankIDCode:              "GBDSC",
			BaseCurrency:            "GBP",
			Bic:                     "NWBKGB42",
			Country:                 ptrStr("GB"),
			CustomerID:              "12345",
			Iban:                    "GB28NWBK40030212764204",
			Name:                    []string{"Jane Doe"},
			NameMatchingStatus:      "opted_out",
			AlternativeNames:        []string{"Sam Holder"},
			JointAccount:            ptrBool(false),
			SecondaryIdentification: "A1B2C3D4",
			PrivateIdentification: &accounts.PrivateIdentification{
				BirthDate:      "2017-07-23",
				BirthCountry:   "GB",
				Identification: "13YH458762",
				Address:        []string{"10 Avenue des Champs"},
				City:           "London",
				Country:        "GB",
				Title:          "Mrs",
				FirstName:      "Jane",
				LastName:       "Doe",
				DocumentNumber: "123456789",
			},
			OrganisationIdentification: &accounts.OrganisationIdentification{
				Identification: "123654",
				Actors: []accounts.Actor{
					{
						Name:      []string{"Jeff Page"},
						BirthDate: "1970-01-01",
						Residency: "GB",
					},
				},
				Address:            []string{"10 Avenue des Champs"},
				City:               "London",
				Country:            "GB",
				Name:               "Jane Doe Ltd",
				RegistrationNumber: "123456789",
				Representative: &accounts.OrganisationRepresentative{
					Name:      "John Smith",
					BirthDate: "1970-01-01",
					Residency: "GB",
				},
			},
			Status:       ptrAccountStatus(accounts.AccountStatusConfirmed),
			StatusReason: "unspecified",
			UserDefinedData: []accounts.UserDefinedData{
				{
					Key:   "Some account related key",
					Value: "Some account related value",
				},
			},
			ValidationType:      "card",
			ReferenceMask:       "############",
			AcceptanceQualifier: "same_day",
			Relationships: &accounts.Relationships{
				MasterAccount: &accounts.RelationshipMasterAccount{
					Data: []accounts.AccountReference{
						{
							Type: "accounts",
							ID:   "a52d13a4-f435-4c00-cfad-f5e7ac5972df",
						},
					},
				},
				AccountEvents: &accounts.RelationshipAccountEvents{
					Data: []accounts.AccountEventReference{
						{
							Type: "account_events",
							ID:   "c1023677-70ee-417a-9a6a-e211241f1e9c",
						},
						{
							Type: "account_events",
							ID:   "437284fa-62a6-4f1d-893d-2959c9780288",
						},
					},
				},
			},
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

	req := accounts.NewCreateRequest(&account)

	resp, err := client.Create(context.Background(), req)
	assert.NoError(t, err)

	expectedresp := &accounts.ResponseCreateAccount{
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
