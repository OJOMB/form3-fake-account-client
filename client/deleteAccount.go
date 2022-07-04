package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/OJOMB/form3-fake-account-client/accounts"
)

// Delete attempts to remove an existing account version
// https://api-docs.form3.tech/api.html#organisation-accounts-fetch
func (c *Client) Delete(ctx context.Context, accountID string, version uint) error {
	if accountID == "" {
		return newInputError("accountID cannot be empty", nil)
	}

	// send request
	path := fmt.Sprintf("%s/%s?version=%d", basev1AccountsPath, accountID, version)
	resp, err := c.delete(ctx, path)
	if err != nil {
		return err
	}

	// read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return newInternalError("failed to read response body", err)
	}

	defer resp.Body.Close()

	// handle error response
	if resp.StatusCode != http.StatusNoContent {
		if len(respBody) > 0 {
			// here the server has returned an error message in the response body
			var apiError accounts.ApiError
			if err := json.Unmarshal(respBody, &apiError); err != nil {
				return newInternalError("failed to unmarshal response body", err)
			}

			return newApiError(fmt.Sprintf("failed to delete account, status code %d", resp.StatusCode), apiError)
		}

		// since the server has not returned an error message
		// we add a descriptive messageto express what went wrong to the user over and above just the status code
		var failureMessage string
		switch resp.StatusCode {
		case http.StatusNotFound:
			failureMessage = "account not found"
		case http.StatusInternalServerError:
			failureMessage = "server error"
		default:
			failureMessage = "received response with unexpected status code from server"
		}

		return newApiError(fmt.Sprintf("failed to delete account, status code %d: %s", resp.StatusCode, failureMessage), nil)
	}

	return nil
}
