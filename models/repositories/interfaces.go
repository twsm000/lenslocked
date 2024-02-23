package repositories

import (
	"io"

	"github.com/twsm000/lenslocked/models/entities"
)

type User interface {
	// Create possible errors:
	//  - repositories.ErrFailedToCreateUser
	Create(user *entities.User) error
	// FindByEmail possible errors:
	//  - repositories.ErrUserNotFound
	FindByEmail(email entities.Email) (*entities.User, error)
	// UpdatePassword possible errors:
	//  - repositories.ErrFailedToUpdateUserPassword
	UpdatePassword(user *entities.User) error

	io.Closer
}

type Session interface {
	Create(session *entities.Session) error
	FindUserByToken(token entities.SessionToken) (*entities.User, error)
	DeleteByToken(token entities.SessionToken) error

	io.Closer
}

type PasswordReset interface {
	Create(reset *entities.PasswordReset) error
	FindPasswordResetAndUserByToken(token entities.SessionToken) (*entities.PasswordReset, *entities.User, error)
	DeleteByID(id uint64) error

	io.Closer
}
