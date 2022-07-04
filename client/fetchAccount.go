package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/OJOMB/form3-fake-account-client/accounts"
)

// Fetch attempts to get an existing account
// https://api-docs.form3.tech/api.html#organisation-accounts-fetch
func (c *Client) Fetch(ctx context.Context, accountID string) (*accounts.Response, error) {
	if accountID == "" {
		return nil, newInputError("accountID cannot be empty", nil)
	}

	// send GET request to the accounts endpoint
	path := fmt.Sprintf("%s/%s", basev1AccountsPath, accountID)
	resp, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}

	// read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, newInternalError("failed to read response body", err)
	}

	defer resp.Body.Close()

	// handle error response
	if resp.StatusCode != http.StatusOK {
		var apiError accounts.ApiError
		if err := json.Unmarshal(respBody, &apiError); err != nil {
			return nil, newInternalError("failed to unmarshal response body", err)
		}

		return nil, newApiError(fmt.Sprintf("failed to fetch account, status code %d", resp.StatusCode), apiError)
	}

	// handle success response
	var fetchAccountResp accounts.Response
	if err := json.Unmarshal(respBody, &fetchAccountResp); err != nil {
		return nil, newInternalError("failed to unmarshal response body", err)
	}

	return &fetchAccountResp, nil
}
