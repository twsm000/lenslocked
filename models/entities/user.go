package entities

import (
	"time"
)

type User struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt *time.Time
	Email     Email
	Password  Hash
}

// ValidateUser possible errors:
//   - ErrInvalidUser
//   - ErrInvalidUserEmail
//   - ErrInvalidPassword
func ValidateUser(u *User) Error {
	if u == nil {
		return NewError(ErrInvalidUser)
	}

	if u.Email.IsEmpty() {
		return NewClientError("Email cannot be empty", ErrInvalidUserEmail)
	}

	if len(u.Password) == 0 {
		return NewClientError("Password cannot be empty", ErrInvalidPassword)
	}

	return nil
}

type UserCreatable struct {
	Email    Email
	Password RawPassword
}

// NewCreatableUser possible errors:
//   - ErrFailedToHashPassword
//   - ErrInvalidUser
//   - ErrInvalidUserEmail
//   - ErrInvalidPassword
func NewCreatableUser(input UserCreatable) (*User, Error) {
	user := User{
		Email: input.Email,
	}

	if err := user.Password.GenerateFrom(input.Password); err != nil {
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
