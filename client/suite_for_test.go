package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OJOMB/form3-fake-account-client/accounts"
)

const (
	// testDataAccountDataAllFieldsRequest has every request field populated including deprecated fields
	// missing created_at, modified_at, and version
	testDataAccountDataAllFieldsFormatStr = `{
		"type": "accounts",
		"id": "1dfaf917-c6d6-4e18-b7e7-972e66492976",
		"organisation_id": "caca9817-6936-4da4-96e7-9ce93206070f",
		"attributes": {
			"account_classification": "Personal",
			"account_number": "10000004",
			"bank_id": "400302",
			"bank_id_code": "GBDSC",
			"base_currency": "GBP",
			"bic": "NWBKGB42",
			"country": "GB",
			"customer_id": "12345",
			"iban": "GB28NWBK40030212764204",
			"name": [
				"Jane Doe"
			],
			"name_matching_status": "opted_out",
			"alternative_names": [
				"Sam Holder"
			],
			"joint_account": false,
			"secondary_identification": "A1B2C3D4",
			"status": "confirmed",
			"status_reason": "unspecified",
			"user_defined_data": [
				{
					"key": "Some account related key",
					"value": "Some account related value"
				}
			],
			"validation_type": "card",
			"reference_mask": "############",
			"acceptance_qualifier": "same_day",
			"title": "Mrs",
			"first_name": "Jane",
			"bank_account_name": "Jane Doe",
			"alternative_bank_account_names": [
				"Sam Holder"
			],
			"processing_service": "processing_service_1",
			"user_defined_information": "Some account related value",
			"account_matching_opt_out": false,
			"switched": false
		},
		"created_on": "%s",
		"modified_on": "%s",
		"version": %d
	}`
)

// testDataAccountDataAllFields returns a JSON string representing an account with all fields populated
func getTestDataAccountDataAllFields(createdOn, modifiedOn string, version int64) string {
	return fmt.Sprintf(testDataAccountDataAllFieldsFormatStr, createdOn, modifiedOn, version)
}

// getDummyTime returns an arbitrary static time for testing purposes
func getDummyTime() time.Time {
	return time.Date(2017, 07, 23, 0, 0, 0, 0, time.UTC)
}

// mockRoundTripper is a mockable round tripper for the http client
// transportFunc serves as the RoundTrip function for the mock
type mockRoundTripper struct {
	transportFunc func(req *http.Request) (*http.Response, error)
}

func (mrt *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return mrt.transportFunc(req)
}

func ptrStr(s string) *string {
	return &s
}

func ptrInt64(i int64) *int64 {
	return &i
}

func ptrBool(b bool) *bool {
	return &b
}

func ptrAccountStatus(s accounts.AccountStatus) *accounts.AccountStatus {
	return &s
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

// errReader is intended to help us test - mainly in the corner case of error handling in the case of defective response bodies
type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed to read")
}
