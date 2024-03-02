package services

import (
	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

type User interface {
	// Create possible errors:
	//   - entities.ErrFailedToHashPassword
	//   - entities.ErrInvalidUser
	//   - entities.ErrInvalidUserEmail
	//   - entities.ErrInvalidPassword
	//   - repositories.ErrFailedToCreateUser
	Create(input entities.UserCreatable) (*entities.User, entities.Error)

	// Authenticate possible errors:
	//   - ErrInvalidAuthCredentials {repositories.ErrUserNotFound, entities.ErrInvalidPassword}
	Authenticate(input entities.UserAuthenticable) (*entities.User, entities.Error)
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

// Create possible errors:
//   - entities.ErrFailedToHashPassword
//   - entities.ErrInvalidUser
//   - entities.ErrInvalidUserEmail
//   - entities.ErrInvalidPassword
//   - repositories.ErrFailedToCreateUser
func (us *userService) Create(input entities.UserCreatable) (*entities.User, entities.Error) {
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
//   - ErrInvalidAuthCredentials {repositories.ErrUserNotFound, entities.ErrInvalidPassword}
func (us *userService) Authenticate(input entities.UserAuthenticable) (*entities.User, entities.Error) {
	const invalidCredentialsErrMsg string = "Invalid credentials."
	user, err := us.Repository.FindByEmail(input.Email)
	if err != nil {
		return nil, entities.NewClientError(invalidCredentialsErrMsg, ErrInvalidAuthCredentials, err)
	}

	if err := user.Password.Compare(input.Password); err != nil {
		return nil, entities.NewClientError(invalidCredentialsErrMsg, ErrInvalidAuthCredentials, err)
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
