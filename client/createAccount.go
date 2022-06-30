package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/OJOMB/form3-fake-account-client/accounts"
)

// Create creates a new account
// https://api-docs.form3.tech/api.html#organisation-accounts-create
func (c *Client) Create(ctx context.Context, req *accounts.RequestCreateAccount) (*accounts.ResponseCreateAccount, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		// return nil, fmt.Errorf("internal error - failed to marshal input request: %v", err)
		return nil, newInternalError("failed to marshal input request", err)
	}

	// Create a new POST request to the accounts endpoint
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.baseURL, basev1AccountsPath), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, newInternalError("failed to create http request", err)
	}

	// send request
	resp, err := c.Do(httpReq)
	if err != nil {
		return nil, newInternalError("failed to send http request", err)
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

		return nil, newApiError(fmt.Sprintf("failed to create account, status code: %d", resp.StatusCode), apiError)
	}

	// handle success response
	var createdAccountResp accounts.ResponseCreateAccount
	if err := json.Unmarshal(respBody, &createdAccountResp); err != nil {
		return nil, newInternalError("failed to unmarshal response body", err)
	}

	return &createdAccountResp, nil
}
