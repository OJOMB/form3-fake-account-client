package client

import "fmt"

type clientErrType int

const (
	apiError clientErrType = iota
	internalError

	apiErrorStr      = "api error"
	internalErrorStr = "internal error"
)

var clientErrors = map[clientErrType]string{
	apiError:      apiErrorStr,
	internalError: internalErrorStr,
}

type clientError struct {
	code clientErrType
	msg  string
	err  error
}

func newInternalError(msg string, err error) *clientError {
	return &clientError{code: internalError, msg: msg, err: err}
}

func newApiError(msg string, err error) *clientError {
	return &clientError{code: apiError, msg: msg, err: err}
}

func (cerr *clientError) Error() string {
	errMsg := fmt.Sprintf("%s - %v", clientErrors[cerr.code], cerr.msg)
	if cerr.err != nil {
		errMsg += fmt.Sprintf(": %v", cerr.err.Error())
	}

	return errMsg
}
