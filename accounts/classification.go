package accounts

import (
	"encoding/json"
	"fmt"
)

type AccountClassification int

const (
	AccountClassificationPersonal AccountClassification = iota
	AccountClassificationBusiness

	accountClassificationPersonalStr = "Personal"
	accountClassificationBusinessStr = "Business"
)

func NewAccountClassification(s string) (AccountClassification, error) {
	switch s {
	case accountClassificationPersonalStr:
		return AccountClassificationPersonal, nil
	case accountClassificationBusinessStr:
		return AccountClassificationBusiness, nil
	default:
		return -1, fmt.Errorf("invalid account classification: %s", s)
	}
}

func (c AccountClassification) String() string {
	switch c {
	case AccountClassificationPersonal:
		return accountClassificationPersonalStr
	case AccountClassificationBusiness:
		return accountClassificationBusinessStr
	default:
		return ""
	}
}

func (c AccountClassification) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *AccountClassification) UnmarshalJSON(data []byte) error {
	var statusStr string
	var err error
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}

	*c, err = NewAccountClassification(statusStr)
	return err
}
