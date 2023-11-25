package postgresrepo

import (
	"database/sql"
	"errors"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

const (
	insertSessionQuery = `
		INSERT INTO sessions (created_at, user_id, token)
		VALUES (CURRENT_TIMESTAMP, $1, $2)
		RETURNING id, created_at
	`
	FindUserByTokenQuery = `
		SELECT u.id,
		       u.created_at,
			   u.updated_at,
			   u.email,
			   u.password
          FROM sessions s
		 INNER JOIN users u
		    ON u.id = s.user_id
		WHERE s.token = $1
	`
)

func NewSessionRepository(db *sql.DB) (*sessionRepository, error) {
	insertSessionStmt, err := db.Prepare(insertSessionQuery)
	if err != nil {
		return nil, err
	}

	findUserByTokenStmt, err := db.Prepare(FindUserByTokenQuery)
	if err != nil {
		return nil, err
	}

	return &sessionRepository{
		db:                  db,
		insertSessionStmt:   insertSessionStmt,
		findUserByTokenStmt: findUserByTokenStmt,
	}, nil
}

type sessionRepository struct {
	db                  *sql.DB
	insertSessionStmt   *sql.Stmt
	findUserByTokenStmt *sql.Stmt
}

func (sr *sessionRepository) Create(userID uint64) (*entities.Session, error) {
	session, err := entities.NewCreatableSession(userID)
	if err != nil {
		return nil, err
	}

	row := sr.insertSessionStmt.QueryRow(session.UserID, session.Token.Hash())
	if err := row.Scan(&session.ID, &session.CreatedAt); err != nil {
		return nil, errors.Join(repositories.ErrFailedToCreateUser, err)
	}

	return session, nil
}

func (sr *sessionRepository) FindUserByToken(token entities.SessionToken) (*entities.User, error) {
	row := sr.findUserByTokenStmt.QueryRow(token.Hash())
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
