package repositories

import "errors"

var (
	ErrFailedToCreateUser          = errors.New("failed to create user")
	ErrFailedToCreateSession       = errors.New("failed to create session")
	ErrFailedToDeleteSession       = errors.New("failed to delete session")
	ErrFailedToCreatePasswordReset = errors.New("failed to create password reset")
	ErrFailedToDeletePasswordReset = errors.New("failed to delete password reset")
	ErrUserNotFound                = errors.New("user not found")
)
