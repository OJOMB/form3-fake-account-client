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
			"private_identification": {
				"birth_date": "2017-07-23",
				"birth_country": "GB",
				"identification": "13YH458762",
				"address": [
					"10 Avenue des Champs"
				],
				"city": "London",
				"country": "GB",
				"title": "Mrs",
				"first_name": "Jane",
				"last_name": "Doe",
				"document_number": "123456789"
			},
			"organisation_identification": {
				"identification": "123654",
				"actors": [
					{
						"name": [
							"Jeff Page"
						],
						"birth_date": "1970-01-01",
						"residency": "GB"
					}
				],
				"address": [
					"10 Avenue des Champs"
				],
				"city": "London",
				"country": "GB",
				"name": "Jane Doe Ltd",
				"registration_number": "123456789",
				"representative": {
					"name": "John Smith",
					"birth_date": "1970-01-01",
					"residency": "GB"
				}
			},
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
			"relationships": {
				"master_account": {
					"data": [
						{
							"type": "accounts",
							"id": "a52d13a4-f435-4c00-cfad-f5e7ac5972df"
						}
					]
				},
				"account_events": {
					"data": [
						{
							"type": "account_events",
							"id": "c1023677-70ee-417a-9a6a-e211241f1e9c"
						},
						{
							"type": "account_events",
							"id": "437284fa-62a6-4f1d-893d-2959c9780288"
						}
					]
				}
			},
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
