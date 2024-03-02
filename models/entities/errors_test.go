package entities

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorWithoutErrors(t *testing.T) {
	err := NewError()
	assert.Nil(t, err)
}

func TestNewError(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	err := NewError(err1)
	assert.EqualValues(t, "error 1", err.Error())
	assert.True(t, err.Is(err1))
	assert.False(t, err.Is(err2))
	assert.False(t, err.IsClientErr())
	assert.Empty(t, err.ClientErr())

	err = NewError(err1, err2)
	assert.EqualValues(t, errors.Join(err1, err2).Error(), err.Error())
	assert.True(t, err.Is(err1))
	assert.True(t, err.Is(err2))
	assert.False(t, err.IsClientErr())
	assert.Empty(t, err.ClientErr())
}

func TestNewClientErrorReturnsNilWhen(t *testing.T) {
	testCases := []struct {
		desc string
		err  error
	}{
		{
			desc: "PassDefaultZeroValues",
			err:  NewClientError(""),
		},
		{
			desc: "PassNoErrors",
			err:  NewClientError("This is not an client error"),
		},
		{
			desc: "WithNilErrorArgs",
			err:  NewClientError("This is not an client error too", nil, nil, nil),
		},
		{
			desc: "WithErrorButWithoutClientErrorMessage",
			err:  NewClientError("", errors.New("error")),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Nil(t, tC.err)
		})
	}
}

func TestNewClientError(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error e")
	t.Run("WithSingleClientError", func(t *testing.T) {
		const clientErrorMsg string = "An error has occured"
		joinErrors := errors.Join(err2, err1)
		err := NewClientError(clientErrorMsg, err2, err1)
		assert.Error(t, err)
		assert.EqualValues(t, joinErrors.Error(), err.Error())
		assert.True(t, err.Is(err1))
		assert.True(t, err.Is(err2))
		assert.True(t, err.IsClientErr())
		assert.EqualValues(t, clientErrorMsg, err.ClientErr())
	})

	t.Run("WithTwoSequentialClientErrors", func(t *testing.T) {
		joinErrors := errors.Join(err2, err1)
		joinClientErrors := strings.Join([]string{"Second fail", "First fail"}, "\n")
		err := NewClientError("Second fail", err2, NewClientError("First fail", err1))
		assert.Error(t, err)
		assert.EqualValues(t, joinErrors.Error(), err.Error())
		assert.True(t, err.Is(err1))
		assert.True(t, err.Is(err2))
		assert.True(t, err.IsClientErr())
		assert.EqualValues(t, joinClientErrors, err.ClientErr())
		var ee *entityError
		assert.True(t, err.As(&ee))
	})

	t.Run("WithTwoClientErrosButNotSequentially", func(t *testing.T) {
		joinErrors := errors.Join(err3, err2, err1)
		// err2 breaks the ClientError sequence
		err := NewClientError("Second fail", err3, err2, NewClientError("First fail", err1))
		assert.Error(t, err)
		assert.EqualValues(t, joinErrors.Error(), err.Error())
		assert.True(t, err.Is(err1))
		assert.True(t, err.Is(err2))
		assert.True(t, err.Is(err3))
		assert.True(t, err.IsClientErr())
		assert.EqualValues(t, "Second fail", err.ClientErr())
	})
}
