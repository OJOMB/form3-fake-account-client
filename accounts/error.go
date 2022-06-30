package accounts

type ApiError struct {
	ErrMsg string `json:"error_message"`
}

func (apierr ApiError) Error() string {
	return apierr.ErrMsg
}
