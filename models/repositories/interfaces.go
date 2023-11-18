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
