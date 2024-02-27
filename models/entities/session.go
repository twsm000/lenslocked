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

// NewCreatableSession possible errors:
//   - rand.ErrFailedToGenerateSlice
//   - rand.ErrInvalidSizeUnexpected
//   - ErrTokenSizeBelowMinRequired
func NewCreatableSession(userID uint64, bytesPerToken int) (*Session, Error) {
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
