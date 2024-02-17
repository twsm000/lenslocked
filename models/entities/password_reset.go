package entities

import "time"

type PasswordReset struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt *time.Time
	UserID    uint64
	Token     SessionToken
	ExpiresAt time.Time
}

func NewCreatablePasswordReset(userID uint64, bytesPerToken int) (*PasswordReset, error) {
	var token SessionToken
	err := token.Update(bytesPerToken)
	if err != nil {
		return nil, err
	}
	pr := PasswordReset{
		UserID: userID,
		Token:  token,
	}
	return &pr, nil
}
