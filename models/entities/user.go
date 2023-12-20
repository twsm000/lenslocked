package entities

import (
	"errors"
	"time"
)

var (
	ErrInvalidUser          = errors.New("invalid user")
	ErrInvalidUserEmail     = errors.New("invalid user email")
	ErrInvalidUserPassword  = errors.New("invalid user password")
	ErrFailedToHashPassword = errors.New("failed to hash password")
)

type User struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt *time.Time
	Email     Email
	Password  Hash
}

func ValidateUser(u *User) error {
	if u == nil {
		return ErrInvalidUser
	}

	if u.Email.IsEmpty() {
		return ErrInvalidUserEmail
	}

	if len(u.Password) == 0 {
		return ErrInvalidUserPassword
	}

	return nil
}

type UserCreatable struct {
	Email    Email
	Password RawPassword
}

func NewCreatableUser(input UserCreatable) (*User, error) {
	user := User{
		Email: input.Email,
	}

	if err := user.Password.GenerateFrom(input.Password.AsBytes()); err != nil {
		return nil, err
	}

	if err := ValidateUser(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

type UserAuthenticable struct {
	Email    Email
	Password RawPassword
}
