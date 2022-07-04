package client

import "fmt"

type clientErrType int

const (
	apiError clientErrType = iota
	internalError
	inputError

	apiErrorStr      = "api error"
	internalErrorStr = "internal error"
	inputErrorStr    = "input error"
)

var clientErrors = map[clientErrType]string{
	apiError:      apiErrorStr,
	internalError: internalErrorStr,
	inputError:    inputErrorStr,
}

// clientError is an error type that is used to represent errors that occur in the client
type clientError struct {
	code clientErrType
	msg  string
	err  error
}

// newInternalError is a helper function that constructs a new clientError of code inputError with the given message and error.
func newInternalError(msg string, err error) *clientError {
	return &clientError{code: internalError, msg: msg, err: err}
}

// newApiError is a helper function that constructs a new clientError of code apiError with the given message and error.
func newApiError(msg string, err error) *clientError {
	return &clientError{code: apiError, msg: msg, err: err}
}

// newInputError is a helper function that constructs a new clientError of code inputError with the given message and error.
func newInputError(msg string, err error) *clientError {
	return &clientError{code: inputError, msg: msg, err: err}
}

func (cerr *clientError) Error() string {
	errMsg := fmt.Sprintf("%s - %v", clientErrors[cerr.code], cerr.msg)
	if cerr.err != nil {
		errMsg += fmt.Sprintf(": %v", cerr.err.Error())
	}

	return errMsg
}
