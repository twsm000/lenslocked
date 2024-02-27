package repositories

import "errors"

var (
	ErrDuplicateUserEmailNotAllowed = errors.New("duplicate user email not allowed")
	ErrFailedToCreateUser           = errors.New("failed to create user")
	ErrFailedToUpdateUserPassword   = errors.New("failed to update user password")
	ErrFailedToCreateSession        = errors.New("failed to create session")
	ErrFailedToDeleteSession        = errors.New("failed to delete session")
	ErrFailedToCreatePasswordReset  = errors.New("failed to create password reset")
	ErrFailedToDeletePasswordReset  = errors.New("failed to delete password reset")
	ErrFixedTokenSizeRequired       = errors.New("fixed token size required")
	ErrUserNotFound                 = errors.New("user not found")
)
