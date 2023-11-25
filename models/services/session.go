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

func (ss Session) FindUserByToken(token entities.SessionToken) (*entities.User, error) {
	return ss.Repository.FindUserByToken(token)
}
