package repositories

import (
	"io"

	"github.com/twsm000/lenslocked/models/entities"
)

type User interface {
	// Create possible errors:
	//  - ErrFailedToCreateUser {ErrDuplicateUserEmailNotAllowed}
	Create(user *entities.User) entities.Error
	// FindByEmail possible errors:
	//  - ErrUserNotFound
	FindByEmail(email entities.Email) (*entities.User, entities.Error)
	// UpdatePassword possible errors:
	//  - ErrFailedToUpdateUserPassword
	UpdatePassword(user *entities.User) error

	io.Closer
}

type Session interface {
	// Create possible errors:
	//   - ErrFailedToCreateSession {ErrFixedTokenSizeRequired, ErrUserNotFound}
	Create(session *entities.Session) entities.Error
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
