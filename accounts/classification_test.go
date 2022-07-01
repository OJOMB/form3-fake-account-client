package accounts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountClassification(t *testing.T) {
	testCases := []struct {
		name                          string
		input                         string
		expectedAccountClassification AccountClassification
		expectedError                 error
	}{
		{
			name:                          "Personal",
			input:                         "Personal",
			expectedAccountClassification: AccountClassificationPersonal,
			expectedError:                 nil,
		},
		{
			name:                          "Business",
			input:                         "Business",
			expectedAccountClassification: AccountClassificationBusiness,
			expectedError:                 nil,
		},
		{
			name:                          "invalid",
			input:                         "invalid",
			expectedAccountClassification: -1,
			expectedError:                 fmt.Errorf("invalid account classification: invalid"),
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx+1, tc.name), func(t *testing.T) {
			accountClassification, err := NewAccountClassification(tc.input)
			assert.Equal(t, tc.expectedAccountClassification, accountClassification)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAccountClassificationString(t *testing.T) {
	testCases := []struct {
		name           string
		input          AccountClassification
		expectedString string
	}{
		{
			name:           "Personal",
			input:          AccountClassificationPersonal,
			expectedString: "Personal",
		},
		{
			name:           "Business",
			input:          AccountClassificationBusiness,
			expectedString: "Business",
		},
		{
			name:           "invalid",
			input:          -1,
			expectedString: "",
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx+1, tc.name), func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.input.String())
		})
	}
}

func TestAccountClassificationUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name                          string
		input                         string
		expectedAccountClassification AccountClassification
		expectedErrorMsg              string
	}{
		{
			name:                          "Personal",
			input:                         `"Personal"`,
			expectedAccountClassification: AccountClassificationPersonal,
		},
		{
			name:                          "Business",
			input:                         `"Business"`,
			expectedAccountClassification: AccountClassificationBusiness,
		},
		{
			name:                          "invalid",
			input:                         `"invalid"`,
			expectedAccountClassification: -1,
			expectedErrorMsg:              "invalid account classification: invalid",
		},
		{
			name:                          "invalid JSON",
			input:                         `invalid"`,
			expectedAccountClassification: 0,
			expectedErrorMsg:              "invalid character 'i' looking for beginning of value",
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx+1, tc.name), func(t *testing.T) {
			var accountClass AccountClassification
			err := accountClass.UnmarshalJSON([]byte(tc.input))
			assert.Equal(t, tc.expectedAccountClassification, accountClass)
			if tc.expectedErrorMsg != "" {
				assert.Equal(t, tc.expectedErrorMsg, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAccountClassificationMarshalJSON(t *testing.T) {
	testCases := []struct {
		name           string
		input          AccountClassification
		expectedString string
	}{
		{
			name:           "Personal",
			input:          AccountClassificationPersonal,
			expectedString: `"Personal"`,
		},
		{
			name:           "Business",
			input:          AccountClassificationBusiness,
			expectedString: `"Business"`,
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
