package repositories

import "github.com/twsm000/lenslocked/models/entities"

type User interface {
	// Create possible errors:
	//  - repositories.ErrFailedToCreateUser
	Create(user *entities.User) error
	// FindByEmail possible errors:
	//  - repositories.ErrUserNotFound
	FindByEmail(email entities.Email) (*entities.User, error)
}

type Session interface {
	Create(session *entities.Session) error
	FindUserByToken(token entities.SessionToken) (*entities.User, error)
}
