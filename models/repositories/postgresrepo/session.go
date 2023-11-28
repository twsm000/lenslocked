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
		ON CONFLICT (user_id)
		DO UPDATE SET token = EXCLUDED.token, updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
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

func NewSessionRepository(db *sql.DB) (repositories.Session, error) {
	insertUpdateSessionStmt, err := db.Prepare(insertSessionQuery)
	if err != nil {
		return nil, err
	}

	findUserByTokenStmt, err := db.Prepare(FindUserByTokenQuery)
	if err != nil {
		return nil, err
	}

	return &sessionRepository{
		db:                      db,
		insertUpdateSessionStmt: insertUpdateSessionStmt,
		findUserByTokenStmt:     findUserByTokenStmt,
	}, nil
}

type sessionRepository struct {
	db                      *sql.DB
	insertUpdateSessionStmt *sql.Stmt
	findUserByTokenStmt     *sql.Stmt
}

func (sr *sessionRepository) Create(session *entities.Session) error {
	row := sr.insertUpdateSessionStmt.QueryRow(session.UserID, session.Token.Hash())
	if err := row.Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt); err != nil {
		return errors.Join(repositories.ErrFailedToCreateSession, err)
	}

	return nil
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
