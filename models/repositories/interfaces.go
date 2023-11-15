package repositories

import "github.com/twsm000/lenslocked/models/entities"

type UserRepository interface {
	Create(input *entities.UserCreatable) (*entities.User, error)
}
