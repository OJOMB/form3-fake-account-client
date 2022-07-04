package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInternalError_constructsCorrectly(t *testing.T) {
	err := newInternalError("test", nil)
	assert.Equal(t, &clientError{code: internalError, msg: "test"}, err)

	originalErr := fmt.Errorf("test1 suffix")
	err = newInternalError("test1", originalErr)
	assert.Equal(t, &clientError{code: internalError, msg: "test1", err: originalErr}, err)
}

func TestNewApiError_constructsCorrectly(t *testing.T) {
	err := newApiError("test", nil)
	assert.Equal(t, &clientError{code: apiError, msg: "test"}, err)

	originalErr := fmt.Errorf("test1 suffix")
	err = newApiError("test1", fmt.Errorf("test1 suffix"))
	assert.Equal(t, &clientError{code: apiError, msg: "test1", err: originalErr}, err)
}

func TestNewInputError_constructsCorrectly(t *testing.T) {
	err := newInputError("test", nil)
	assert.Equal(t, &clientError{code: inputError, msg: "test"}, err)

	originalErr := fmt.Errorf("test1 suffix")
	err = newInputError("test1", fmt.Errorf("test1 suffix"))
	assert.Equal(t, &clientError{code: inputError, msg: "test1", err: originalErr}, err)
}

func TestInternalError_ErrorMsgFormatsCorrectly(t *testing.T) {
	err := newInternalError("test", nil)
	assert.Equal(t, "internal error - test", err.Error())

	err = newInternalError("test1", fmt.Errorf("test1 suffix"))
	assert.Equal(t, "internal error - test1: test1 suffix", err.Error())
}

func TestApiError_ErrorMsgFormatsCorrectly(t *testing.T) {
	err := newApiError("test", nil)
	assert.Equal(t, "api error - test", err.Error())

	err = newApiError("test1", fmt.Errorf("test1 suffix"))
	assert.Equal(t, "api error - test1: test1 suffix", err.Error())
}

func TestInputError_ErrorMsgFormatsCorrectly(t *testing.T) {
	err := newInputError("test", nil)
	assert.Equal(t, "input error - test", err.Error())

	err = newInputError("test1", fmt.Errorf("test1 suffix"))
	assert.Equal(t, "input error - test1: test1 suffix", err.Error())
}
