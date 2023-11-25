package entities

import (
	"fmt"

	"github.com/twsm000/lenslocked/pkg/http/session"
)

type SessionToken struct {
	Hash string
	// Value is only available when Update or UpdateDefault is called.
	// Normally this occurs when create a new session or updating an existing one.
	Value string
}

func (st *SessionToken) UpdateDefault() error {
	return st.Update(32)
}

func (st *SessionToken) Update(size int) error {
	token, err := session.Token(size)
	if err != nil {
		return err
	}
	st.Hash = "" // TODO: hash
	st.Value = token
	return nil
}

func (st *SessionToken) String() string {
	return hiddenHash
}

func (st *SessionToken) Scan(value any) error {
	if value == nil {
		st.Value = ""
		st.Hash = ""
		return nil
	}

	hash, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid scan type: %T", value)
	}
	st.Hash = hash
	st.Value = ""
	return nil
}
