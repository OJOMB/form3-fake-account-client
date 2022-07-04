package accounts

type (
	Request struct {
		Data *AccountData `json:"data"`
	}

	Response struct {
		Data  *AccountData `json:"data"`
		Links *Links       `json:"links"`
	}

	Links struct {
		Self string `json:"self"`
	}
)

func NewRequest(data AccountData) *Request {
	return &Request{Data: &data}
}
