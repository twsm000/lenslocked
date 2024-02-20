package services

import (
	"errors"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

type User interface {
	Create(input entities.UserCreatable) (*entities.User, error)
	Authenticate(input entities.UserAuthenticable) (*entities.User, error)
	UpdatePassword(user *entities.User, rawPassword entities.RawPassword) error
}

func NewUser(repo repositories.User) User {
	return &userService{
		Repository: repo,
	}
}

type userService struct {
	Repository repositories.User
}

func (us *userService) Create(input entities.UserCreatable) (*entities.User, error) {
	user, err := entities.NewCreatableUser(input)
	if err != nil {
		return nil, err
	}

	if err = us.Repository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Authenticate possible errors:
//
// services.ErrInvalidAuthCredentials:
//   - entities.ErrInvalidUserPassword
//   - repositories.ErrUserNotFound
func (us *userService) Authenticate(input entities.UserAuthenticable) (*entities.User, error) {
	user, err := us.Repository.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.Join(ErrInvalidAuthCredentials, err)
	}

	if err := user.Password.Compare(input.Password); err != nil {
		return nil, errors.Join(ErrInvalidAuthCredentials, err)
	}

	return user, nil
}

func (us *userService) UpdatePassword(user *entities.User, rawPassword entities.RawPassword) error {
	if err := user.Password.GenerateFrom(rawPassword); err != nil {
		return err
	}

	if err := entities.ValidateUser(user); err != nil {
		return err
	}

	return us.Repository.UpdatePassword(user)
}
