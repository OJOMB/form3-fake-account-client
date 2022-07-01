package accounts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountStatus(t *testing.T) {
	testCases := []struct {
		name                  string
		input                 string
		expectedAccountStatus AccountStatus
		expectedError         error
	}{
		{
			name:                  "confirmed",
			input:                 "confirmed",
			expectedAccountStatus: AccountStatusConfirmed,
			expectedError:         nil,
		},
		{
			name:                  "pending",
			input:                 "pending",
			expectedAccountStatus: AccountStatusPending,
			expectedError:         nil,
		},
		{
			name:                  "cancelled",
			input:                 "cancelled",
			expectedAccountStatus: AccountStatusCancelled,
			expectedError:         nil,
		},
		{
			name:                  "failed",
			input:                 "failed",
			expectedAccountStatus: AccountStatusFailed,
			expectedError:         nil,
		},
		{
			name:                  "invalid",
			input:                 "invalid",
			expectedAccountStatus: -1,
			expectedError:         fmt.Errorf("invalid account status: invalid"),
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx+1, tc.name), func(t *testing.T) {
			accountStatus, err := NewAccountStatus(tc.input)
			assert.Equal(t, tc.expectedAccountStatus, accountStatus)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAccountStatusString(t *testing.T) {
	testCases := []struct {
		name           string
		input          AccountStatus
		expectedString string
	}{
		{
			name:           "confirmed",
			input:          AccountStatusConfirmed,
			expectedString: "confirmed",
		},
		{
			name:           "pending",
			input:          AccountStatusPending,
			expectedString: "pending",
		},
		{
			name:           "cancelled",
			input:          AccountStatusCancelled,
			expectedString: "cancelled",
		},
		{
			name:           "failed",
			input:          AccountStatusFailed,
			expectedString: "failed",
		},
		{
			name:           "invalid",
			input:          -1,
			expectedString: "",
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx+1, tc.name), func(t *testing.T) {
			string := tc.input.String()
			assert.Equal(t, tc.expectedString, string)
		})
	}
}

func TestAccountStatusUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name                  string
		input                 string
		expectedAccountStatus AccountStatus
		expectedErrorMsg      string
	}{
		{
			name:                  "confirmed",
			input:                 `"confirmed"`,
			expectedAccountStatus: AccountStatusConfirmed,
		},
		{
			name:                  "pending",
			input:                 `"pending"`,
			expectedAccountStatus: AccountStatusPending,
		},
		{
			name:                  "cancelled",
			input:                 `"cancelled"`,
			expectedAccountStatus: AccountStatusCancelled,
		},
		{
			name:                  "failed",
			input:                 `"failed"`,
			expectedAccountStatus: AccountStatusFailed,
		},
		{
			name:                  "invalid value",
			input:                 `"invalid"`,
			expectedAccountStatus: -1,
			expectedErrorMsg:      "invalid account status: invalid",
		},
		{
			name:                  "invalid JSON",
			input:                 `invalid"`,
			expectedAccountStatus: 0,
			expectedErrorMsg:      "invalid character 'i' looking for beginning of value",
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx+1, tc.name), func(t *testing.T) {
			var accountStatus AccountStatus
			err := accountStatus.UnmarshalJSON([]byte(tc.input))
			assert.Equal(t, tc.expectedAccountStatus, accountStatus)
			if tc.expectedErrorMsg != "" {
				assert.Equal(t, tc.expectedErrorMsg, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAccountStatusMarshalJSON(t *testing.T) {
	testCases := []struct {
		name           string
		input          AccountStatus
		expectedString string
	}{
		{
			name:           "confirmed",
			input:          AccountStatusConfirmed,
			expectedString: `"confirmed"`,
		},
		{
			name:           "pending",
			input:          AccountStatusPending,
			expectedString: `"pending"`,
		},
		{
			name:           "cancelled",
			input:          AccountStatusCancelled,
			expectedString: `"cancelled"`,
		},
		{
			name:           "failed",
			input:          AccountStatusFailed,
			expectedString: `"failed"`,
		},
		{
			name:           "invalid",
			input:          -1,
			expectedString: `""`,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx+1, tc.name), func(t *testing.T) {
			bytes, err := tc.input.MarshalJSON()
			assert.Equal(t, tc.expectedString, string(bytes))
			assert.Equal(t, nil, err)
		})
	}
}
