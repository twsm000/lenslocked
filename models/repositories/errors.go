package repositories

import "errors"

var (
	ErrFailedToCreateUser = errors.New("failed to create user")
	ErrUserNotFound       = errors.New("user not found")
)
