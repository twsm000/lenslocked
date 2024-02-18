package postgresrepo

import (
	"database/sql"
	"errors"
	"log"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

const (
	insertPasswordResetQuery = `
		INSERT INTO password_resets (created_at, user_id, token, expires_at)
		VALUES (CURRENT_TIMESTAMP, $1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET token = EXCLUDED.token
		             ,updated_at = CURRENT_TIMESTAMP
					 ,expires_at = EXCLUDED.expires_at
		RETURNING id, created_at, updated_at
	`
	findUserByPasswordResetTokenQuery = `
		SELECT u.id,
               u.created_at,
               u.updated_at,
               u.email,
               u.password
          FROM password_resets r
		 INNER JOIN users u
		    ON u.id = r.user_id
		WHERE r.token = $1
	`

	deleteByPasswordResetTokenQuery = `
		DELETE FROM password_resets
		 WHERE token = $1
	`
)

func NewPasswordResetRepository(db *sql.DB, logErr, logInfo, logWarn *log.Logger) (repositories.PasswordReset, error) {
	insertUpdateStmt, err := db.Prepare(insertPasswordResetQuery)
	if err != nil {
		return nil, err
	}

	findUserByTokenStmt, err := db.Prepare(findUserByPasswordResetTokenQuery)
	if err != nil {
		return nil, err
	}

	deleteByTokenStmt, err := db.Prepare(deleteByPasswordResetTokenQuery)
	if err != nil {
		return nil, err
	}

	return &passwordResetRepository{
		db:                  db,
		logErr:              logErr,
		logInfo:             logInfo,
		logWarn:             logWarn,
		insertUpdateStmt:    insertUpdateStmt,
		findUserByTokenStmt: findUserByTokenStmt,
		deleteByTokenStmt:   deleteByTokenStmt,
	}, nil
}

type passwordResetRepository struct {
	db                  *sql.DB
	logErr              *log.Logger
	logInfo             *log.Logger
	logWarn             *log.Logger
	insertUpdateStmt    *sql.Stmt
	findUserByTokenStmt *sql.Stmt
	deleteByTokenStmt   *sql.Stmt
}

func (sr *passwordResetRepository) Close() error {
	return errors.Join(
		sr.deleteByTokenStmt.Close(),
		sr.findUserByTokenStmt.Close(),
		sr.insertUpdateStmt.Close(),
	)
}

func (sr *passwordResetRepository) Create(reset *entities.PasswordReset) error {
	row := sr.insertUpdateStmt.QueryRow(reset.UserID, reset.Token.Hash(), reset.ExpiresAt)
	if err := row.Scan(&reset.ID, &reset.CreatedAt, &reset.UpdatedAt); err != nil {
		return errors.Join(repositories.ErrFailedToCreatePasswordReset, err)
	}

	return nil
}

func (sr *passwordResetRepository) FindUserByToken(token entities.SessionToken) (*entities.User, error) {
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

func (sr *passwordResetRepository) DeleteByToken(token entities.SessionToken) error {
	result, err := sr.deleteByTokenStmt.Exec(token.Hash())
	if err != nil {
		return errors.Join(repositories.ErrFailedToDeletePasswordReset, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	switch rowsAffected {
	case 0:
		sr.logWarn.Println("Try to delete password reset, but not found:", token.Value())
	case 1:
		sr.logInfo.Println("PasswordReset deleted successfully:", token.Value())
	default:
		sr.logErr.Printf("Failed to delete password reset: %q, rows affected: %d", token.Value(), rowsAffected)
	}
	return nil
}
