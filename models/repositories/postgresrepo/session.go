package postgresrepo

import (
	"database/sql"
	"errors"
	"log"
	"strings"

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
	findUserBySessionTokenQuery = `
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

	deleteBySessionTokenQuery = `
		DELETE FROM sessions
		 WHERE token = $1
	`
)

func NewSessionRepository(db *sql.DB, logErr, logInfo, logWarn *log.Logger) (repositories.Session, error) {
	insertUpdateSessionStmt, err := db.Prepare(insertSessionQuery)
	if err != nil {
		return nil, err
	}

	findUserByTokenStmt, err := db.Prepare(findUserBySessionTokenQuery)
	if err != nil {
		return nil, err
	}

	deleteByTokenStmt, err := db.Prepare(deleteBySessionTokenQuery)
	if err != nil {
		return nil, err
	}

	return &sessionRepository{
		db:                      db,
		logErr:                  logErr,
		logInfo:                 logInfo,
		logWarn:                 logWarn,
		insertUpdateSessionStmt: insertUpdateSessionStmt,
		findUserByTokenStmt:     findUserByTokenStmt,
		deleteByTokenStmt:       deleteByTokenStmt,
	}, nil
}

type sessionRepository struct {
	db                      *sql.DB
	logErr                  *log.Logger
	logInfo                 *log.Logger
	logWarn                 *log.Logger
	insertUpdateSessionStmt *sql.Stmt
	findUserByTokenStmt     *sql.Stmt
	deleteByTokenStmt       *sql.Stmt
}

func (sr *sessionRepository) Close() error {
	return errors.Join(
		sr.deleteByTokenStmt.Close(),
		sr.findUserByTokenStmt.Close(),
		sr.insertUpdateSessionStmt.Close(),
	)
}

// Create possible errors:
//   - ErrFailedToCreateSession {ErrFixedTokenSizeRequired, ErrUserNotFound}
func (sr *sessionRepository) Create(session *entities.Session) entities.Error {
	row := sr.insertUpdateSessionStmt.QueryRow(session.UserID, session.Token.Hash())
	if err := row.Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt); err != nil {
		if strings.Contains(err.Error(), "sessions_token_check") {
			return entities.NewError(
				repositories.ErrFailedToCreateSession,
				repositories.ErrFixedTokenSizeRequired,
				err,
			)
		} else if strings.Contains(err.Error(), "sessions_user_id_fkey") {
			return entities.NewClientError(
				"User not found",
				repositories.ErrFailedToCreateSession,
				repositories.ErrUserNotFound,
				err,
			)
		}
		return entities.NewError(repositories.ErrFailedToCreateSession, err)
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

func (sr *sessionRepository) DeleteByToken(token entities.SessionToken) error {
	result, err := sr.deleteByTokenStmt.Exec(token.Hash())
	if err != nil {
		return errors.Join(repositories.ErrFailedToDeleteSession, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	switch rowsAffected {
	case 0:
		sr.logWarn.Println("Try to delete session, but not found:", token.Value())
	case 1:
		sr.logInfo.Println("Session deleted successfully:", token.Value())
	default:
		sr.logErr.Printf("Failed to delete session: %q, rows affected: %d", token.Value(), rowsAffected)
	}
	return nil
}
