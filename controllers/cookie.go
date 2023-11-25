package controllers

import (
	"net/http"

	"github.com/twsm000/lenslocked/models/entities"
)

const (
	CookieSession = "session"
)

func createSessionCookie(session *entities.Session) *http.Cookie {
	return &http.Cookie{
		Name:     CookieSession,
		Value:    session.Token.Value(),
		Path:     "/",
		HttpOnly: true,
	}
}
