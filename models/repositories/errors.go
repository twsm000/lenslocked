package repositories

import "errors"

var (
	ErrFailedToCreateUser    = errors.New("failed to create user")
	ErrFailedToCreateSession = errors.New("failed to create session")
	ErrFailedToDeleteSession = errors.New("failed to delete session")
	ErrUserNotFound          = errors.New("user not found")
)
