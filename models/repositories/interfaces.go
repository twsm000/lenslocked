package repositories

import "github.com/twsm000/lenslocked/models/entities"

type User interface {
	// Create possible errors:
	//  - repositories.ErrFailedToCreateUser
	Create(input entities.UserCreatable) (*entities.User, error)
	// FindByEmail possible errors:
	//  - repositories.ErrUserNotFound
	FindByEmail(email entities.Email) (*entities.User, error)
}

type Session interface {
	Create(userID uint64) (*entities.Session, error)
	FindUserByToken(token entities.SessionToken) (*entities.User, error)
}
