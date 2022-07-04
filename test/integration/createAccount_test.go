package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/OJOMB/form3-fake-account-client/accounts"
	"github.com/OJOMB/form3-fake-account-client/client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount_SuccessPath(t *testing.T) {
	c, err := client.NewClient(testBaseURL, nil)
	assert.NoError(t, err)

	accountID, err := uuid.NewRandom()
	assert.NoError(t, err)

	account := accounts.AccountData{
		ID:             accountID.String(),
		OrganisationID: "600b4bf3-4cae-4e1c-b382-968f86fc7489",
		Type:           "accounts",
		Attributes: &accounts.AccountAttributes{
			AccountClassification:   accounts.AccountClassificationPersonal,
			AccountMatchingOptOut:   false,
			AccountNumber:           "10000004",
			AlternativeNames:        []string{"Sam Holder"},
			BankID:                  "400302",
			BankIDCode:              "GBDSC",
			BaseCurrency:            "GBP",
			Bic:                     "NWBKGB42",
			Country:                 ptrStr("GB"),
			Iban:                    "GB28NWBK40030212764204",
			JointAccount:            ptrBool(false),
			Name:                    []string{"Jane Doe"},
			SecondaryIdentification: "A1B2C3D4",
			Status:                  ptrAccountStatus(accounts.AccountStatusConfirmed),
			Switched:                ptrBool(false),
		},
	}

	beforeTest := time.Now()

	resp, err := c.Create(context.Background(), account)
	assert.NoError(t, err)

	afterTest := time.Now()

	// quick check to see that CreatedOn and ModifiedOn times fall in valid time ranges
	assert.True(t, resp.Data.CreatedOn.After(beforeTest) && resp.Data.CreatedOn.Before(afterTest))
	assert.True(t, resp.Data.ModifiedOn.After(beforeTest) && resp.Data.ModifiedOn.Before(afterTest))

	expectedResp := &accounts.Response{
		Data: &account,
		Links: &accounts.Links{
			Self: fmt.Sprintf("/v1/organisation/accounts/%s", accountID),
		},
	}
	expectedResp.Data.Version = ptrInt64(0)
	expectedResp.Data.CreatedOn = ptrTime(getDummyTime())
	expectedResp.Data.ModifiedOn = ptrTime(getDummyTime())

	// overwrite the time values in the response to match the dummy value
	resp.Data.CreatedOn = ptrTime(getDummyTime())
	resp.Data.ModifiedOn = ptrTime(getDummyTime())

	assert.Equal(t, expectedResp, resp)
}

// TestCreateAccount_WithExistingID_FailurePath checks for expected errors when creating an account with an ID that already exists
func TestCreateAccount_WithExistingID_FailurePath(t *testing.T) {
	c, err := client.NewClient(testBaseURL, nil)
	assert.NoError(t, err)

	accountID, err := uuid.NewRandom()
	assert.NoError(t, err)

	account := accounts.AccountData{
		ID:             accountID.String(),
		OrganisationID: "600b4bf3-4cae-4e1c-b382-968f86fc7489",
		Type:           "accounts",
		Attributes: &accounts.AccountAttributes{
			Country: ptrStr("GB"),
			Name:    []string{"Jane Doe"},
		},
	}

	resp1, err := c.Create(context.Background(), account)
	assert.NoError(t, err)
	assert.NotNil(t, resp1)

	// we know the previous account created successfully, now let's try to recreate it and hope it fails
	resp2, err := c.Create(context.Background(), account)
	assert.Error(t, err)
	assert.Nil(t, resp2)

	// check the error message
	assert.Equal(t, "api error - failed to create account, status code 409: Account cannot be created as it violates a duplicate constraint", err.Error())
}

// TestCreateAccount_WithMissingRequiredFields_FailurePath checks for expected errors when required fields are omitted
func TestCreateAccount_WithMissingRequiredFields_FailurePath(t *testing.T) {
	const missingDataErrorMsgFormatNest3 = "api error - failed to create account, status code 400: validation failure list:\nvalidation failure list:\nvalidation failure list:\n%s in body is required"
	const missingDataErrorMsgFormatNest2 = "api error - failed to create account, status code 400: validation failure list:\nvalidation failure list:\n%s in body is required"

	c, err := client.NewClient(testBaseURL, nil)
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		accountData      accounts.AccountData
		expectedErrorMsg string
	}{
		{
			name: "missing country",
			accountData: accounts.AccountData{
				OrganisationID: "600b4bf3-4cae-4e1c-b382-968f86fc7489",
				Type:           "accounts",
				Attributes:     &accounts.AccountAttributes{Name: []string{"Jane Doe"}},
			},
			expectedErrorMsg: fmt.Sprintf(missingDataErrorMsgFormatNest3, "country"),
		},
		{
			name: "missing name",
			accountData: accounts.AccountData{
				OrganisationID: "600b4bf3-4cae-4e1c-b382-968f86fc7489",
				Type:           "accounts",
				Attributes:     &accounts.AccountAttributes{Country: ptrStr("GB")},
			},
			expectedErrorMsg: fmt.Sprintf(missingDataErrorMsgFormatNest3, "name"),
		},
		{
			name: "missing organisation ID",
			accountData: accounts.AccountData{
				Type: "accounts",
				Attributes: &accounts.AccountAttributes{
					Country: ptrStr("GB"),
					Name:    []string{"Jane Doe"},
				},
			},
			expectedErrorMsg: fmt.Sprintf(missingDataErrorMsgFormatNest2, "organisation_id"),
		},
		{
			name: "missing type",
			accountData: accounts.AccountData{
				OrganisationID: "600b4bf3-4cae-4e1c-b382-968f86fc7489",
				Attributes:     &accounts.AccountAttributes{Country: ptrStr("GB"), Name: []string{"Jane Doe"}},
			},
			expectedErrorMsg: fmt.Sprintf(missingDataErrorMsgFormatNest2, "type"),
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx+1, tc.name), func(t *testing.T) {
			accountID, err := uuid.NewRandom()
			assert.NoError(t, err)

			tc.accountData.ID = accountID.String()

			resp, err := c.Create(context.Background(), tc.accountData)
			assert.Error(t, err)
			assert.Nil(t, resp)

			// check the error message
			assert.Equal(t, tc.expectedErrorMsg, err.Error())
		})
	}
}
