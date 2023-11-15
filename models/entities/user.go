package entities

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint64
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	Email        Email
	PasswordHash UserPasswordHash
}

type Email struct {
	email string
}

func (e *Email) Set(email string) {
	e.email = strings.ToLower(email)
}

func (e Email) String() string {
	return e.email
}

type UserPasswordHash []byte

const hiddenPassword string = "********"

func (up UserPasswordHash) AsBytes() []byte {
	return up
}

func (up UserPasswordHash) String() string {
	return hiddenPassword
}

var ErrUserPasswordHash = errors.New("failed to hash user password")

func (up *UserPasswordHash) GenerateFrom(password []byte) error {
	passwordHashed, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return errors.Join(ErrUserPasswordHash, err)
	}
	*up = passwordHashed
	return nil
}

type UserCreatable struct {
	Email    string
	Password []byte
}
