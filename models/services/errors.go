package services

import "errors"

var (
	ErrInvalidAuthCredentials    = errors.New("invalid authentication credentials")
	ErrPasswordResetTokenExpired = errors.New("password reset token expired")
)
