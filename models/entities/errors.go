package entities

import (
	"errors"
	"strings"
)

var (
	ErrFailedToHashPassword = errors.New("failed to hash password")
	ErrInvalidTokenSize     = errors.New("invalid token size")
	ErrInvalidUser          = errors.New("invalid user")
	ErrInvalidUserEmail     = errors.New("invalid user email")
	ErrInvalidPassword      = errors.New("invalid password")
)

// Error is an interface to complement the error interface
// and ensure that an error ocurred in the domain model is
// for the client or for the application.
type Error interface {
	error
	ClientErr() string
	IsClientErr() bool
	As(target any) bool
	Is(target error) bool
}

// NewError join multiples errors for the application level and returns.
// Nil will be returned if none or all args are nil
// The errors are joined by errors.Join
func NewError(errs ...error) Error {
	err := errors.Join(errs...)
	if err == nil {
		return nil
	}

	return &entityError{
		err: errors.Join(errs...),
	}
}

// NewClientError join multiples errors for the client level and returns.
// Nil will be returnd if the clientErr string has default zero value or
// errors.Join return nil.
func NewClientError(clientErr string, errs ...error) Error {
	err := errors.Join(errs...)
	if err == nil || clientErr == "" {
		return nil
	}

	return &entityError{
		clientErr: clientErr,
		err:       err,
	}
}

// entityError implements Error inteface
type entityError struct {
	clientErr string
	err       error
}

// Error returns the joined errors
func (e *entityError) Error() string {
	return e.err.Error()
}

// ClientErr returns an empty string when IsClientErr is false,
// otherwise, returns all client errors joined by new line character.
// This multiple client errors will be joined when the next errors
// implements then ClientError inteface. If a non ClientError is found,
// returns the result of the already collected strings.
func (e *entityError) ClientErr() string {
	if !e.IsClientErr() {
		return ""
	}

	result := []string{e.clientErr}
	if ue, ok := e.err.(interface{ Unwrap() []error }); ok {
		for _, err := range ue.Unwrap()[1:] {
			if cerr, ok := err.(ClientError); ok {
				result = append(result, cerr.ClientErr())
			} else {
				break
			}
		}
	}
	return strings.Join(result, "\n")
}

// IsClientErr returns true when Error the current error is an client error
func (e *entityError) IsClientErr() bool {
	return e.clientErr != ""
}

// Is wrap the errors.Is and apply the target error to the joined errors
func (e *entityError) Is(target error) bool {
	return errors.Is(e.err, target)
}

// As wrap then errors.As and apply the target error to the joined errors
func (e *entityError) As(target any) bool {
	return errors.As(e.err, target)
}

// ClientError represents an error condition that can be displayed
// to clients
type ClientError interface {
	ClientErr() string
}
