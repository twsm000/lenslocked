package services

import (
	"errors"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

type User struct {
	Repository repositories.User
}

func (us User) Create(input entities.UserCreatable) (*entities.User, error) {
	return us.Repository.Create(input)
}

// Authenticate possible errors:
//
// services.ErrInvalidAuthCredentials:
//   - entities.ErrInvalidUserPassword
//   - repositories.ErrUserNotFound
func (us User) Authenticate(input entities.UserAuthenticable) (*entities.User, error) {
	user, err := us.Repository.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.Join(ErrInvalidAuthCredentials, err)
	}
	if err := user.Password.Compare(input.Password); err != nil {
		return nil, errors.Join(ErrInvalidAuthCredentials, err)
	}
	return user, nil
}
