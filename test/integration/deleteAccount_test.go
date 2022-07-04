package integration

import (
	"context"
	"testing"

	"github.com/OJOMB/form3-fake-account-client/accounts"
	"github.com/OJOMB/form3-fake-account-client/client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestDeleteAccount_SuccessPath tests the success path of the DeleteAccount function
// this is done by creatiung an account and then deleting it. Verifying that we get a successful response.
func TestDeleteAccount_SuccessPath(t *testing.T) {
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

	// first need to successfully create the account
	createResp, err := c.Create(context.Background(), account)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	// now we can attempt to delete the account version 0
	err = c.Delete(context.Background(), accountID.String(), 0)
	assert.NoError(t, err)
}

func TestDeleteAccount_DeleteNonExistentAccount_Failure(t *testing.T) {
	c, err := client.NewClient(testBaseURL, nil)
	assert.NoError(t, err)

	accountID, err := uuid.NewRandom()
	assert.NoError(t, err)

	// now we attempt to delete the account version 0 which should fail since it doesn't exist
	err = c.Delete(context.Background(), accountID.String(), 0)
	assert.Equal(t, "api error - failed to delete account, status code 404: account not found", err.Error())
}

func TestDeleteAccount_IncorrectVersion_SuccessPath(t *testing.T) {
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

	// first need to successfully create the account
	createResp, err := c.Create(context.Background(), account)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	// now we can attempt to delete the account version 1 knowing the account is at version 0
	err = c.Delete(context.Background(), accountID.String(), 1)
	assert.Equal(t, "api error - failed to delete account, status code 409: invalid version", err.Error())
}

func TestDeleteAccount_InvalidUUID_SuccessPath(t *testing.T) {
	c, err := client.NewClient(testBaseURL, nil)
	assert.NoError(t, err)

	// now we can attempt to delete an account with an invalid uuid
	err = c.Delete(context.Background(), "invalid-uuid", 1)
	assert.Equal(t, "api error - failed to delete account, status code 400: id is not a valid uuid", err.Error())
}
