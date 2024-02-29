package entities

import "errors"

var (
	ErrFailedToHashPassword = errors.New("failed to hash password")
	ErrInvalidTokenSize     = errors.New("invalid token size")
	ErrInvalidUser          = errors.New("invalid user")
	ErrInvalidUserEmail     = errors.New("invalid user email")
	ErrInvalidUserPassword  = errors.New("invalid user password")
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

// NewError return join multiples errors for the application level
func NewError(errs ...error) Error {
	return &entityError{
		err: errors.Join(errs...),
	}
}

// NewClientError join multiples errors for the client level
func NewClientError(clientErr string, errs ...error) Error {
	return &entityError{
		clientErr: clientErr,
		err:       errors.Join(errs...),
	}
}

// entityError implements Error inteface
type entityError struct {
	clientErr string
	err       error
}

func (e *entityError) Error() string {
	return e.err.Error()
}

func (e *entityError) ClientErr() string {
	return e.clientErr
}

func (e *entityError) IsClientErr() bool {
	return e.clientErr != ""
}

func (e *entityError) Is(target error) bool {
	return errors.Is(e.err, target)
}

func (e *entityError) As(target any) bool {
	return errors.As(e.err, target)
}

type ClientError interface {
	ClientErr() string
}
