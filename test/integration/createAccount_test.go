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
	c := client.NewClient(testBaseURL, nil)

	accountID, err := uuid.NewRandom()
	assert.NoError(t, err)

	account := &accounts.AccountData{
		ID:             accountID.String(),
		OrganisationID: "600b4bf3-4cae-4e1c-b382-968f86fc7489",
		Type:           "accounts",
		Attributes: &accounts.AccountAttributes{
			AccountClassification:   ptrStr("Personal"),
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

	req := accounts.NewCreateRequest(account)

	beforeTest := time.Now()

	resp, err := c.Create(context.Background(), req)
	assert.NoError(t, err)

	afterTest := time.Now()

	// quick check to see that CreatedOn and ModifiedOn are in valid time ranges
	assert.True(t, resp.Data.CreatedOn.After(beforeTest) && resp.Data.CreatedOn.Before(afterTest))
	assert.True(t, resp.Data.ModifiedOn.After(beforeTest) && resp.Data.ModifiedOn.Before(afterTest))

	expectedResp := &accounts.ResponseCreateAccount{
		Data: account,
		Links: &accounts.Links{
			Self: fmt.Sprintf("/v1/organisation/accounts/%s", accountID),
		},
	}
	expectedResp.Data.Version = ptrInt64(0)
	expectedResp.Data.CreatedOn = ptrTime(getDummyTime())
	expectedResp.Data.ModifiedOn = ptrTime(getDummyTime())

	// fix the time values in the response to match the dummy value
	resp.Data.CreatedOn = ptrTime(getDummyTime())
	resp.Data.ModifiedOn = ptrTime(getDummyTime())

	assert.Equal(t, expectedResp, resp)
}
