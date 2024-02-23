package postgresrepo

import (
	"database/sql"
	"errors"
	"log"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

const (
	queryInsertPasswordReset = `
		INSERT INTO password_resets (created_at, user_id, token, expires_at)
		VALUES (CURRENT_TIMESTAMP, $1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET token = EXCLUDED.token
		             ,updated_at = CURRENT_TIMESTAMP
					 ,expires_at = EXCLUDED.expires_at
		RETURNING id, created_at, updated_at
	`
	queryFindPasswordResetAndUserByToken = `
		SELECT pr.id,
               pr.created_at,
               pr.updated_at,
               pr.user_id,
               pr.token,
               pr.expires_at,
               u.id,
               u.created_at,
               u.updated_at,
               u.email,
               u.password
          FROM password_resets pr
		 INNER JOIN users u
		    ON u.id = pr.user_id
		WHERE pr.token = $1
	`

	QueryDeletePasswordResetByID = `
		DELETE FROM password_resets
		 WHERE id = $1
	`
)

func NewPasswordResetRepository(db *sql.DB, logErr, logInfo, logWarn *log.Logger) (repositories.PasswordReset, error) {
	insertUpdateStmt, err := db.Prepare(queryInsertPasswordReset)
	if err != nil {
		return nil, err
	}

	findUserByTokenStmt, err := db.Prepare(queryFindPasswordResetAndUserByToken)
	if err != nil {
		return nil, err
	}

	deleteByTokenStmt, err := db.Prepare(QueryDeletePasswordResetByID)
	if err != nil {
		return nil, err
	}

	return &passwordResetRepository{
		db:                             db,
		logErr:                         logErr,
		logInfo:                        logInfo,
		logWarn:                        logWarn,
		insertUpdateStmt:               insertUpdateStmt,
		findPasswordAndUserByTokenStmt: findUserByTokenStmt,
		deleteByTokenStmt:              deleteByTokenStmt,
	}, nil
}

type passwordResetRepository struct {
	db                             *sql.DB
	logErr                         *log.Logger
	logInfo                        *log.Logger
	logWarn                        *log.Logger
	insertUpdateStmt               *sql.Stmt
	findPasswordAndUserByTokenStmt *sql.Stmt
	deleteByTokenStmt              *sql.Stmt
}

func (sr *passwordResetRepository) Close() error {
	return errors.Join(
		sr.deleteByTokenStmt.Close(),
		sr.findPasswordAndUserByTokenStmt.Close(),
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

func (sr *passwordResetRepository) FindPasswordResetAndUserByToken(
	token entities.SessionToken) (*entities.PasswordReset, *entities.User, error) {
	/*******************************************************************************/
	row := sr.findPasswordAndUserByTokenStmt.QueryRow(token.Hash())
	var passwordReset entities.PasswordReset
	var user entities.User
	err := row.Scan(
		&passwordReset.ID,
		&passwordReset.CreatedAt,
		&passwordReset.UpdatedAt,
		&passwordReset.UserID,
		&passwordReset.Token,
		&passwordReset.ExpiresAt,
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, nil, errors.Join(repositories.ErrUserNotFound, err)
	}
	return &passwordReset, &user, nil
}

func (sr *passwordResetRepository) DeleteByID(id uint64) error {
	result, err := sr.deleteByTokenStmt.Exec(id)
	if err != nil {
		return errors.Join(repositories.ErrFailedToDeletePasswordReset, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	switch rowsAffected {
	case 0:
		sr.logWarn.Println("Try to delete password reset, but not found:", id)
	case 1:
		sr.logInfo.Println("PasswordReset deleted successfully:", id)
	default:
		sr.logErr.Printf("Failed to delete password reset: %d, rows affected: %d", id, rowsAffected)
	}
	return nil
}
