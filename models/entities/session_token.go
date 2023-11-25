package entities

import "github.com/twsm000/lenslocked/pkg/http/session"

type SessionToken string

func (st *SessionToken) UpdateDefault() error {
	return st.Update(32)
}

func (st *SessionToken) Update(size int) error {
	t, err := session.Token(size)
	if err != nil {
		return err
	}
	*st = SessionToken(t)
	return nil
}
