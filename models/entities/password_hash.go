package entities

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	hiddenHash string = "********"
)

type RawPassword string

func (rp *RawPassword) Set(password string) {
	*rp = RawPassword(password)
}

func (rp RawPassword) AsBytes() []byte {
	return []byte(rp)
}

type Hash []byte

func (h Hash) AsBytes() []byte {
	return h
}

func (h Hash) String() string {
	return hiddenHash
}

// Compare possible errors:
//   - ErrInvalidPassword
func (h Hash) Compare(rawPassword RawPassword) Error {
	if err := bcrypt.CompareHashAndPassword(h, rawPassword.AsBytes()); err != nil {
		return NewClientError("The given password does not match with your current password.", ErrInvalidPassword, err)
	}
	return nil
}

// GenerateFrom possible errors:
//   - ErrFailedToHashPassword
func (h *Hash) GenerateFrom(rawPassword RawPassword) Error {
	passwordHashed, err := bcrypt.GenerateFromPassword(rawPassword.AsBytes(), bcrypt.DefaultCost)
	if err != nil {
		return NewError(ErrFailedToHashPassword, err)
	}
	*h = passwordHashed
	return nil
}

func (h *Hash) Scan(value any) error {
	if value == nil {
		*h = nil
		return nil
	}

	hash, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid scan type: %T", value)
	}
	*h = []byte(hash)
	return nil
}
