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
	FindUserByToken(token entities.SessionToken) (*entities.User, error)
	DeleteByToken(token entities.SessionToken) error

	io.Closer
}
