package services

import (
	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

type UserService struct {
	Repository repositories.UserRepository
}

func (us UserService) Create(input *entities.UserCreatable) (*entities.User, error) {
	// todo: validate data
	return us.Repository.Create(input)
}
