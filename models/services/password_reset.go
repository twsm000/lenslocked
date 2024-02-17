package services

import (
	"time"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

const (
	DefaultPasswordResetDuration = 1 * time.Hour
)

type PasswordReset interface {
	Create(email string) (*entities.PasswordReset, error)
	Consume(token string) (*entities.User, error)
}

func NewPasswordReset(bytesPerToken int, repo repositories.PasswordReset) PasswordReset {
	if bytesPerToken < entities.MinBytesPerToken {
		bytesPerToken = entities.MinBytesPerToken
	}

	return PasswordResetService{
		BytesPerToken: bytesPerToken,
		Repository:    repo,
	}
}

type PasswordResetService struct {
	BytesPerToken int
	Repository    repositories.PasswordReset
	// Duration is the amount of time that a PasswordReset is valid for
	Duration time.Duration
}

func (prs PasswordResetService) Create(email string) (*entities.PasswordReset, error) {
	var user entities.User // TODO: Find user by email
	passwordReset, err := entities.NewCreatablePasswordReset(user.ID, prs.BytesPerToken)
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
}
