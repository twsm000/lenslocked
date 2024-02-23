package services

import (
	"log"
	"time"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

type PasswordReset interface {
	Create(email entities.Email) (*entities.PasswordReset, error)
	Consume(token entities.SessionToken) (*entities.User, error)
}

func NewPasswordReset(
	bytesPerToken int,
	duration time.Duration,
	repo repositories.PasswordReset,
	userRepo repositories.User,
	logError *log.Logger) PasswordReset {
	/***************************************************/
	if bytesPerToken < entities.MinBytesPerToken {
		bytesPerToken = entities.MinBytesPerToken
	}

	return PasswordResetService{
		BytesPerToken:  bytesPerToken,
		Repository:     repo,
		Duration:       duration,
		UserRepository: userRepo,
		logError:       logError,
	}
}

type PasswordResetService struct {
	BytesPerToken int

	// Duration is the amount of time that a PasswordReset is valid for
	Duration       time.Duration
	Repository     repositories.PasswordReset
	UserRepository repositories.User

	// logs
	logError *log.Logger
}

func (prs PasswordResetService) Create(email entities.Email) (*entities.PasswordReset, error) {
	user, err := prs.UserRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	passwordReset, err := entities.NewCreatablePasswordReset(
		user.ID,
		prs.BytesPerToken,
		entities.NewPasswordResetTimeout(time.Now(), prs.Duration),
	)
	if err != nil {
		return nil, err
	}

	if err := prs.Repository.Create(passwordReset); err != nil {
		return nil, err
	}

	return passwordReset, err
}

func (prs PasswordResetService) Consume(token entities.SessionToken) (*entities.User, error) {
	passwordReset, user, err := prs.Repository.FindPasswordResetAndUserByToken(token)
	if err != nil {
		return nil, err
	}
	defer prs.Repository.DeleteByID(passwordReset.ID) // error ignored because its not useful

	now := time.Now()
	if passwordReset.ExpiresAt.Before(now) {
		prs.logError.Printf("\nPasswordReset.ExpiresAt: %v - Now: %v\n", passwordReset.ExpiresAt, now)
		return nil, ErrPasswordResetTokenExpired
	}

	return user, nil
}
