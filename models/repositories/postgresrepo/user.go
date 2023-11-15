package postgresrepo

import (
	"database/sql"
	"errors"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *sql.DB
}

func (ur *userRepository) Create(input *entities.UserCreatable) (*entities.User, error) {
	return nil, errors.New("not implemented")
}
