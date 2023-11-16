package repositories

import "github.com/twsm000/lenslocked/models/entities"

type User interface {
	Create(input *entities.User) (*entities.User, error)
}
