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
)

func NewSessionRepository(db *sql.DB) (*sessionRepository, error) {
	insertSessionStmt, err := db.Prepare(insertSessionQuery)
	if err != nil {
		return nil, err
	}

	return &sessionRepository{
		db:                db,
		insertSessionStmt: insertSessionStmt,
	}, nil
}

type sessionRepository struct {
	db                *sql.DB
	insertSessionStmt *sql.Stmt
}

func (sr *sessionRepository) Create(userID uint64) (*entities.Session, error) {
	session, err := entities.NewCreatableSession(userID)
	if err != nil {
		return nil, err
	}

	row := sr.insertSessionStmt.QueryRow(session.UserID, session.Token.Value)
	if err := row.Scan(&session.ID, &session.CreatedAt); err != nil {
		return nil, errors.Join(repositories.ErrFailedToCreateUser, err)
	}

	return session, nil
}
