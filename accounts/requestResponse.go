package accounts

type (
	RequestCreateAccount struct {
		Data *AccountData `json:"data"`
	}

	ResponseCreateAccount struct {
		Data  *AccountData `json:"data"`
		Links *Links       `json:"links"`
	}

	Links struct {
		Self string `json:"self"`
	}
)

func NewCreateRequest(data *AccountData) *RequestCreateAccount {
	return &RequestCreateAccount{
		Data: data,
	}
}
