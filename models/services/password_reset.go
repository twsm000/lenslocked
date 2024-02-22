package services

import (
	"errors"
	"time"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

type PasswordReset interface {
	Create(email entities.Email) (*entities.PasswordReset, error)
	Consume(token string) (*entities.User, error)
}

func NewPasswordReset(bytesPerToken int, duration time.Duration, repo repositories.PasswordReset,
	userRepo repositories.User) PasswordReset {
	if bytesPerToken < entities.MinBytesPerToken {
		bytesPerToken = entities.MinBytesPerToken
	}

	return PasswordResetService{
		BytesPerToken:  bytesPerToken,
		Repository:     repo,
		Duration:       duration,
		UserRepository: userRepo,
	}
}

type PasswordResetService struct {
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for
	Duration       time.Duration
	Repository     repositories.PasswordReset
	UserRepository repositories.User
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

func (prs PasswordResetService) Consume(token string) (*entities.User, error) {
	// TODO: implement prs.Consume
	// var stoken entities.SessionToken
	// err := stoken.SetFromHex(token)
	// if err != nil {
	// 	return nil, err
	// }
	// return prs.Repository.FindUserByToken(stoken)
	return nil, errors.New("not implemented")
}
