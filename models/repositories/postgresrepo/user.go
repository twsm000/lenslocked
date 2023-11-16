package postgresrepo

import (
	"database/sql"
	"errors"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

const (
	insertUser = `
		INSERT INTO users (created_at, email, password)
		VALUES (CURRENT_TIMESTAMP, $1, $2)
		RETURNING id, created_at
	`
)

func NewUserRepository(db *sql.DB) (repositories.User, error) {
	insertUserStmt, err := db.Prepare(insertUser)
	if err != nil {
		return nil, err
	}
	return &userRepository{
		db:             db,
		insertUserStmt: insertUserStmt,
	}, nil
}

type userRepository struct {
	db             *sql.DB
	insertUserStmt *sql.Stmt
}

func (ur *userRepository) Create(input *entities.User) (*entities.User, error) {
	if input == nil {
		return nil, errors.New("invalid user")
	}

	row := ur.insertUserStmt.QueryRow(input.Email, input.PasswordHash.AsBytes())
	err := row.Scan(
		&input.ID,
		&input.CreatedAt,
	)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create user"), err)
	}
	return input, nil
}
