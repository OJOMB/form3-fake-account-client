package accounts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountNameMatchingStatus(t *testing.T) {
	testCases := []struct {
		name                              string
		input                             string
		expectedAccountNameMatchingStatus AccountNameMatchingStatus
		expectedError                     error
	}{
		{
			name:                              "supported",
			input:                             "supported",
			expectedAccountNameMatchingStatus: AccountNameMatchingStatusSupported,
			expectedError:                     nil,
		},
		{
			name:                              "switched",
			input:                             "switched",
			expectedAccountNameMatchingStatus: AccountNameMatchingStatusSwitched,
			expectedError:                     nil,
		},
		{
			name:                              "opted out",
			input:                             "opted_out",
			expectedAccountNameMatchingStatus: AccountNameMatchingStatusOptedOut,
			expectedError:                     nil,
		},
		{
			name:                              "not supported",
			input:                             "not_supported",
			expectedAccountNameMatchingStatus: AccountNameMatchingStatusNotsupported,
			expectedError:                     nil,
		},
		{
			name:                              "invalid",
			input:                             "invalid",
			expectedAccountNameMatchingStatus: -1,
			expectedError:                     fmt.Errorf("invalid account name matching status: invalid"),
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx+1, tc.name), func(t *testing.T) {
			accountNameMatchingStatus, err := NewAccountNameMatchingStatus(tc.input)
			assert.Equal(t, tc.expectedAccountNameMatchingStatus, accountNameMatchingStatus)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAccountNameMatchingStatusString(t *testing.T) {
	testCases := []struct {
		name           string
		input          AccountNameMatchingStatus
		expectedString string
	}{
		{
			name:           "supported",
			input:          AccountNameMatchingStatusSupported,
			expectedString: "supported",
		},
		{
			name:           "switched",
			input:          AccountNameMatchingStatusSwitched,
			expectedString: "switched",
		},
		{
			name:           "opted out",
			input:          AccountNameMatchingStatusOptedOut,
			expectedString: "opted_out",
		},
		{
			name:           "not supported",
			input:          AccountNameMatchingStatusNotsupported,
			expectedString: "not_supported",
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

func TestAccountNameMatchingStatusUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name                              string
		input                             string
		expectedAccountNameMatchingStatus AccountNameMatchingStatus
		expectedErrorMsg                  string
	}{
		{
			name:                              "supported",
			input:                             `"supported"`,
			expectedAccountNameMatchingStatus: AccountNameMatchingStatusSupported,
		},
		{
			name:                              "switched",
			input:                             `"switched"`,
			expectedAccountNameMatchingStatus: AccountNameMatchingStatusSwitched,
		},
		{
			name:                              "opted out",
			input:                             `"opted_out"`,
			expectedAccountNameMatchingStatus: AccountNameMatchingStatusOptedOut,
		},
		{
			name:                              "not supported",
			input:                             `"not_supported"`,
			expectedAccountNameMatchingStatus: AccountNameMatchingStatusNotsupported,
		},
		{
			name:                              "invalid value",
			input:                             `"invalid"`,
			expectedAccountNameMatchingStatus: -1,
			expectedErrorMsg:                  "invalid account name matching status: invalid",
		},
		{
			name:                              "invalid JSON",
			input:                             `invalid"`,
			expectedAccountNameMatchingStatus: 0,
			expectedErrorMsg:                  "invalid character 'i' looking for beginning of value",
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx+1, tc.name), func(t *testing.T) {
			var accountNameMatchingStatus AccountNameMatchingStatus
			err := accountNameMatchingStatus.UnmarshalJSON([]byte(tc.input))
			assert.Equal(t, tc.expectedAccountNameMatchingStatus, accountNameMatchingStatus)
			if tc.expectedErrorMsg != "" {
				assert.Equal(t, tc.expectedErrorMsg, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAccountNameMatchingStatusMarshalJSON(t *testing.T) {
	testCases := []struct {
		name           string
		input          AccountNameMatchingStatus
		expectedString string
	}{
		{
			name:           "supported",
			input:          AccountNameMatchingStatusSupported,
			expectedString: `"supported"`,
		},
		{
			name:           "switched",
			input:          AccountNameMatchingStatusSwitched,
			expectedString: `"switched"`,
		},
		{
			name:           "opted out",
			input:          AccountNameMatchingStatusOptedOut,
			expectedString: `"opted_out"`,
		},
		{
			name:           "not supported",
			input:          AccountNameMatchingStatusNotsupported,
			expectedString: `"not_supported"`,
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
