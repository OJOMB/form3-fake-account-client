package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/OJOMB/form3-fake-account-client/accounts"
)

// Create attempts to create a new account
// https://api-docs.form3.tech/api.html#organisation-accounts-create
func (c *Client) Create(ctx context.Context, account accounts.AccountData) (*accounts.Response, error) {
	// create request
	req := accounts.NewRequest(account)
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, newInternalError("failed to marshal input request", err)
	}

	// send request
	resp, err := c.post(ctx, basev1AccountsPath, reqBody)
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
	if resp.StatusCode != http.StatusCreated {
		var apiError accounts.ApiError
		if err := json.Unmarshal(respBody, &apiError); err != nil {
			return nil, newInternalError("failed to unmarshal response body", err)
		}

		return nil, newApiError(fmt.Sprintf("failed to create account, status code %d", resp.StatusCode), apiError)
	}

	// handle success response
	var createdAccountResp accounts.Response
	if err := json.Unmarshal(respBody, &createdAccountResp); err != nil {
		return nil, newInternalError("failed to unmarshal response body", err)
	}

	return &createdAccountResp, nil
}
