package services

import (
	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

type Session struct {
	Repository repositories.Session
}

func (ss Session) Create(userID uint64) (*entities.Session, error) {
	return ss.Repository.Create(userID)
}

func (ss Session) FindUserByToken(token string) (*entities.User, error) {
	var stoken entities.SessionToken
	err := stoken.Set(token)
	if err != nil {
		return nil, err
	}
	return ss.Repository.FindUserByToken(stoken)
}
