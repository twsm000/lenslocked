package services

import (
	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

type Session interface {
	Create(userID uint64) (*entities.Session, error)
	FindUserByToken(token string) (*entities.User, error)
}

func NewSession(bytesPerToken int, repo repositories.Session) Session {
	if bytesPerToken < entities.MinBytesPerSessionToken {
		bytesPerToken = entities.MinBytesPerSessionToken
	}
	return sessionService{
		BytesPerToken: bytesPerToken,
		Repository:    repo,
	}
}

type sessionService struct {
	BytesPerToken int
	Repository    repositories.Session
}

func (ss sessionService) Create(userID uint64) (*entities.Session, error) {
	session, err := entities.NewCreatableSession(userID, ss.BytesPerToken)
	if err != nil {
		return nil, err
	}
	if err := ss.Repository.Create(session); err != nil {
		return nil, err
	}
	return session, err
}

func (ss sessionService) FindUserByToken(token string) (*entities.User, error) {
	var stoken entities.SessionToken
	err := stoken.SetFromHex(token)
	if err != nil {
		return nil, err
	}
	return ss.Repository.FindUserByToken(stoken)
}

func (ss sessionService) DeleteByToken(token string) error {
	var stoken entities.SessionToken
	err := stoken.SetFromHex(token)
	if err != nil {
		return err
	}
	return ss.Repository.DeleteByToken(stoken)
}
