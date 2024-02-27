package postgresrepo

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

const (
	insertUserQuery = `
		INSERT INTO users (created_at, email, password)
		VALUES (CURRENT_TIMESTAMP, $1, $2)
		RETURNING id, created_at
	`

	findUserByEmailQuery = `
		SELECT id,
               created_at,
			   updated_at,
			   email,
			   password
		  FROM users
		 where email = $1
	`

	updateUserPasswordQuery = `
		UPDATE users
		   SET password = $2
		 WHERE id = $1
	`
)

func NewUserRepository(db *sql.DB) (repositories.User, error) {
	insertUserStmt, err := db.Prepare(insertUserQuery)
	if err != nil {
		return nil, err
	}

	findUserByEmailStmt, err := db.Prepare(findUserByEmailQuery)
	if err != nil {
		return nil, err
	}

	updateUserPasswordStmt, err := db.Prepare(updateUserPasswordQuery)
	if err != nil {
		return nil, err
	}

	return &userRepository{
		db:                     db,
		insertUserStmt:         insertUserStmt,
		findUserByEmailStmt:    findUserByEmailStmt,
		updateUserPasswordStmt: updateUserPasswordStmt,
	}, nil
}

type userRepository struct {
	db                     *sql.DB
	insertUserStmt         *sql.Stmt
	findUserByEmailStmt    *sql.Stmt
	updateUserPasswordStmt *sql.Stmt
}

func (ur *userRepository) Close() error {
	return errors.Join(
		ur.findUserByEmailStmt.Close(),
		ur.insertUserStmt.Close(),
		ur.updateUserPasswordStmt.Close(),
	)
}

// Create possible errors:
//   - ErrFailedToCreateUser {ErrDuplicateUserEmailNotAllowed}
func (ur *userRepository) Create(user *entities.User) entities.Error {
	row := ur.insertUserStmt.QueryRow(user.Email, user.Password.AsBytes())
	if err := row.Scan(&user.ID, &user.CreatedAt); err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return entities.NewClientError(
				"This email is already used by an user",
				repositories.ErrFailedToCreateUser,
				repositories.ErrDuplicateUserEmailNotAllowed,
				err,
			)
		}
		return entities.NewError(repositories.ErrFailedToCreateUser, err)
	}

	return nil
}

func (ur *userRepository) FindByEmail(email entities.Email) (*entities.User, error) {
	row := ur.findUserByEmailStmt.QueryRow(email)
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

func (ur *userRepository) UpdatePassword(user *entities.User) error {
	if _, err := ur.updateUserPasswordStmt.Exec(user.ID, user.Password.AsBytes()); err != nil {
		return errors.Join(repositories.ErrFailedToUpdateUserPassword, err)
	}
	return nil
}
