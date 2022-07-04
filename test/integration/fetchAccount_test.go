package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/OJOMB/form3-fake-account-client/accounts"
	"github.com/OJOMB/form3-fake-account-client/client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFetchAccount_SuccessPath(t *testing.T) {
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

	// first need to successfully create the account
	createResp, err := c.Create(context.Background(), account)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	// now we can attempt to fetch the account
	fetchResp, err := c.Fetch(context.Background(), accountID.String())
	assert.NoError(t, err)
	assert.NotNil(t, fetchResp)

	// check that the fetched account matches the created account
	expectedAccount := &accounts.Response{
		Data: &account,
		Links: &accounts.Links{
			Self: fmt.Sprintf("/v1/organisation/accounts/%s", accountID.String()),
		},
	}

	// first we have to fix these variables which are missing in the create request/unpredictable
	expectedAccount.Data.CreatedOn = ptrTime(getDummyTime())
	expectedAccount.Data.ModifiedOn = ptrTime(getDummyTime())
	expectedAccount.Data.Version = ptrInt64(0)
	fetchResp.Data.CreatedOn = ptrTime(getDummyTime())
	fetchResp.Data.ModifiedOn = ptrTime(getDummyTime())

	assert.EqualValues(t, expectedAccount, fetchResp)
}

func TestFetchAccount_NonExistentAccount_FailurePath(t *testing.T) {
	c, err := client.NewClient(testBaseURL, nil)
	assert.NoError(t, err)

	// liklehood of collision essentially zero so just use random uuid
	accountID, err := uuid.NewRandom()
	assert.NoError(t, err)

	fetchResp, err := c.Fetch(context.Background(), accountID.String())
	assert.Error(t, err)
	assert.Nil(t, fetchResp)

	// check for expected error
	expectedErrMsg := fmt.Sprintf("api error - failed to fetch account, status code 404: record %s does not exist", accountID.String())
	assert.Equal(t, expectedErrMsg, err.Error())
}

func TestFetchAccount_WithInvalidUUID_FailurePath(t *testing.T) {
	c, err := client.NewClient(testBaseURL, nil)
	assert.NoError(t, err)

	accountID := "not-a-uuid"
	fetchResp, err := c.Fetch(context.Background(), accountID)
	assert.Error(t, err)
	assert.Nil(t, fetchResp)

	// check for expected error
	expectedErrMsg := "api error - failed to fetch account, status code 400: id is not a valid uuid"
	assert.Equal(t, expectedErrMsg, err.Error())
}
