package session

import (
	"errors"

	"github.com/twsm000/lenslocked/pkg/crypto/rand"
)

var (
	ErrFailedToGenerateToken = errors.New("failed to generate token")
)

func Token(size int) (string, error) {
	str, err := rand.String(size)
	if err != nil {
		return "", errors.Join(ErrFailedToGenerateToken, err)
	}
	return str, nil
}
