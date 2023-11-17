package services

import (
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
