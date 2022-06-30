package integration

import (
	"time"

	"github.com/OJOMB/form3-fake-account-client/accounts"
)

const (
	testBaseURL = "http://localhost:8080"
)

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

// getDummyTime returns an arbitrary static time for testing purposes
func getDummyTime() time.Time {
	return time.Date(2017, 07, 23, 0, 0, 0, 0, time.UTC)
}
