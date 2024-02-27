package services

import (
	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

type Session interface {
	// Create possible errors:
	//   - rand.ErrFailedToGenerateSlice
	//   - rand.ErrInvalidSizeUnexpected
	//   - ErrTokenSizeBelowMinRequired
	//   - repositories.ErrFailedToCreateSession {ErrFixedTokenSizeRequired, ErrUserNotFound}
	Create(userID uint64) (*entities.Session, entities.Error)
	FindUserByToken(token string) (*entities.User, error)
	DeleteByToken(token string) error
}

func NewSession(bytesPerToken int, repo repositories.Session) Session {
	if bytesPerToken < entities.MinBytesPerToken {
		bytesPerToken = entities.MinBytesPerToken
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

// Create possible errors:
//   - rand.ErrFailedToGenerateSlice
//   - rand.ErrInvalidSizeUnexpected
//   - ErrTokenSizeBelowMinRequired
//   - repositories.ErrFailedToCreateSession {ErrFixedTokenSizeRequired, ErrUserNotFound}
func (ss sessionService) Create(userID uint64) (*entities.Session, entities.Error) {
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
