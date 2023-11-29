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
	findByEmail = `
		SELECT id,
               created_at,
			   updated_at,
			   email,
			   password
		  FROM users
		 where email = $1
	`
)

func NewUserRepository(db *sql.DB) (repositories.User, error) {
	insertUserStmt, err := db.Prepare(insertUser)
	if err != nil {
		return nil, err
	}
	findByEmailStmt, err := db.Prepare(findByEmail)
	if err != nil {
		return nil, err
	}
	return &userRepository{
		db:              db,
		insertUserStmt:  insertUserStmt,
		findByEmailStmt: findByEmailStmt,
	}, nil
}

type userRepository struct {
	db              *sql.DB
	insertUserStmt  *sql.Stmt
	findByEmailStmt *sql.Stmt
}

func (ur *userRepository) Close() error {
	return errors.Join(
		ur.findByEmailStmt.Close(),
		ur.insertUserStmt.Close(),
	)
}

func (ur *userRepository) Create(user *entities.User) error {
	row := ur.insertUserStmt.QueryRow(user.Email, user.Password.AsBytes())
	if err := row.Scan(&user.ID, &user.CreatedAt); err != nil {
		return errors.Join(repositories.ErrFailedToCreateUser, err)
	}

	return nil
}

func (ur *userRepository) FindByEmail(email entities.Email) (*entities.User, error) {
	row := ur.findByEmailStmt.QueryRow(email)
	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, errors.Join(repositories.ErrUserNotFound, err)
	}
	return &user, nil
}
