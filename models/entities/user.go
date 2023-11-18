package entities

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUser         = errors.New("invalid user")
	ErrInvalidUserEmail    = errors.New("invalid user email")
	ErrInvalidUserPassword = errors.New("invalid user password")
)

type User struct {
	ID           uint64
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	Email        Email
	PasswordHash UserPasswordHash
}

func (u *User) Validate() error {
	if u == nil {
		return ErrInvalidUser
	}

	if u.Email.IsEmpty() {
		return ErrInvalidUserEmail
	}

	if len(u.PasswordHash) == 0 {
		return ErrInvalidUserPassword
	}

	return nil
}

type Email struct {
	email string
}

func (e *Email) Set(email string) {
	e.email = strings.ToLower(email)
}

func (e Email) IsEmpty() bool {
	return e.email == ""
}

func (e Email) String() string {
	return e.email
}

func (e *Email) Scan(value any) error {
	if value == nil {
		e.email = ""
		return nil
	}

	email, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid scan type: %T", value)
	}
	e.Set(email)
	return nil
}

const hiddenPassword string = "********"

type UserPasswordHash []byte

func (up UserPasswordHash) AsBytes() []byte {
	return up
}

func (up UserPasswordHash) String() string {
	return hiddenPassword
}

func (up UserPasswordHash) Compare(password []byte) error {
	if err := bcrypt.CompareHashAndPassword(up, password); err != nil {
		return errors.Join(ErrInvalidUserPassword, err)
	}
	return nil
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

type UserAuthenticator struct {
	Email    Email
	Password []byte
}
