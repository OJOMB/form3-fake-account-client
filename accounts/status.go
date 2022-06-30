package accounts

import (
	"encoding/json"
	"fmt"
)

type AccountStatus int

const (
	AccountStatusConfirmed AccountStatus = iota
	AccountStatusPending
	AccountStatusCancelled
	AccountStatusFailed

	accountStatusConfirmedStr = "confirmed"
	accountStatusPendingStr   = "pending"
	accountStatusCancelledStr = "cancelled"
	accountStatusFailedStr    = "failed"
)

func NewAccountStatus(s string) (AccountStatus, error) {
	switch s {
	case accountStatusConfirmedStr:
		return AccountStatusConfirmed, nil
	case accountStatusPendingStr:
		return AccountStatusPending, nil
	case accountStatusCancelledStr:
		return AccountStatusCancelled, nil
	case accountStatusFailedStr:
		return AccountStatusFailed, nil
	default:
		return 0, fmt.Errorf("invalid account status: %s", s)
	}
}

func (s AccountStatus) String() string {
	switch s {
	case AccountStatusConfirmed:
		return accountStatusConfirmedStr
	case AccountStatusPending:
		return accountStatusPendingStr
	case AccountStatusCancelled:
		return accountStatusCancelledStr
	case AccountStatusFailed:
		return accountStatusFailedStr
	default:
		return ""
	}
}

func (s AccountStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *AccountStatus) UnmarshalJSON(data []byte) error {
	var statusStr string
	var err error
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}

	*s, err = NewAccountStatus(statusStr)
	return err
}
