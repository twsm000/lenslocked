package entities

import (
	"time"
)

type Session struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt *time.Time
	UserID    uint64
	Token     SessionToken
}

func NewCreatableSession(userID uint64, bytesPerToken int) (*Session, error) {
	var token SessionToken
	err := token.Update(bytesPerToken)
	if err != nil {
		return nil, err
	}
	session := Session{
		UserID: userID,
		Token:  token,
	}
	return &session, nil
}
