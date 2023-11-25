package entities

import (
	"errors"
	"fmt"

	"github.com/twsm000/lenslocked/pkg/http/session"
)

var (
	ErrEmptyTokenNotAllowed = errors.New("empty token not allowed")
)

const (
	MinBytesPerSessionToken int = 32
)

type SessionToken struct {
	hash  string
	value string
}

func (st *SessionToken) Update(size int) error {
	token, err := session.Token(size)
	if err != nil {
		return err
	}

	return st.Set(token)
}

func (st *SessionToken) Set(token string) error {
	if token == "" {
		return ErrEmptyTokenNotAllowed
	}

	st.value = token
	st.hash = token // TODO: hash the st.value
	return nil
}

func (st SessionToken) Value() string {
	return st.value
}

func (st SessionToken) Hash() string {
	return st.hash
}

func (st *SessionToken) String() string {
	return hiddenHash
}

func (st *SessionToken) Scan(value any) error {
	st.value = ""
	if value == nil {
		st.hash = ""
		return nil
	}

	hash, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid scan type: %T", value)
	}
	st.hash = hash
	return nil
}
