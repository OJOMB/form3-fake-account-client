package accounts

import (
	"encoding/json"
	"fmt"
)

type AccountNameMatchingStatus int

const (
	AccountNameMatchingStatusSupported AccountNameMatchingStatus = iota
	AccountNameMatchingStatusSwitched
	AccountNameMatchingStatusOptedOut
	AccountNameMatchingStatusNotsupported

	accountNameMatchingStatusSupportedStr    = "supported"
	accountNameMatchingStatusSwitchedStr     = "switched"
	accountNameMatchingStatusOptedOutStr     = "opted_out"
	accountNameMatchingStatusNotsupportedStr = "not_supported"
)

func NewAccountNameMatchingStatus(s string) (AccountNameMatchingStatus, error) {
	switch s {
	case accountNameMatchingStatusSupportedStr:
		return AccountNameMatchingStatusSupported, nil
	case accountNameMatchingStatusSwitchedStr:
		return AccountNameMatchingStatusSwitched, nil
	case accountNameMatchingStatusOptedOutStr:
		return AccountNameMatchingStatusOptedOut, nil
	case accountNameMatchingStatusNotsupportedStr:
		return AccountNameMatchingStatusNotsupported, nil
	default:
		return -1, fmt.Errorf("invalid account name matching status: %s", s)
	}
}

func (nms AccountNameMatchingStatus) String() string {
	switch nms {
	case AccountNameMatchingStatusSupported:
		return accountNameMatchingStatusSupportedStr
	case AccountNameMatchingStatusSwitched:
		return accountNameMatchingStatusSwitchedStr
	case AccountNameMatchingStatusOptedOut:
		return accountNameMatchingStatusOptedOutStr
	case AccountNameMatchingStatusNotsupported:
		return accountNameMatchingStatusNotsupportedStr
	default:
		return ""
	}
}

func (nms AccountNameMatchingStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(nms.String())
}

func (nms *AccountNameMatchingStatus) UnmarshalJSON(data []byte) error {
	var statusStr string
	var err error
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}

	*nms, err = NewAccountNameMatchingStatus(statusStr)
	return err
}
