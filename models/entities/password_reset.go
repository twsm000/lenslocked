package entities

import "time"

const (
	DefaultPasswordResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt *time.Time
	UserID    uint64
	Token     SessionToken
	ExpiresAt time.Time
}

// NewCreatablePasswordReset possible errors:
//   - rand.ErrFailedToGenerateSlice
//   - rand.ErrInvalidSizeUnexpected
//   - ErrTokenSizeBelowMinRequired
func NewCreatablePasswordReset(userID uint64, bytesPerToken int, timeout time.Time) (*PasswordReset, Error) {
	var token SessionToken
	err := token.Update(bytesPerToken)
	if err != nil {
		return nil, err
	}

	pr := PasswordReset{
		UserID:    userID,
		Token:     token,
		ExpiresAt: timeout,
	}
	return &pr, nil
}

func NewPasswordResetTimeout(start time.Time, duration time.Duration) time.Time {
	if duration == 0 {
		duration = DefaultPasswordResetDuration
	}

	return start.Add(duration)
}
