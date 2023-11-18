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
	var user entities.User
	user.Email.Set(input.Email)
	if err := user.PasswordHash.GenerateFrom(input.Password); err != nil {
		return nil, err
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return us.Repository.Create(&user)
}

func (us User) Authenticate(input entities.UserAuthenticator) (*entities.User, error) {
	user, err := us.Repository.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if err := user.PasswordHash.Compare(input.Password); err != nil {
		return nil, errors.Join(ErrInvalidAuthCredentials, err)
	}
	return user, nil
}
