package controllers

import (
	"net/http"

	"github.com/twsm000/lenslocked/models/entities"
)

const (
	CookieSession = "session"
)

func createSessionCookie(session *entities.Session) *http.Cookie {
	return createCookie(CookieSession, session.Token.Value())
}

func createCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
}
func deleteCookie(name string) *http.Cookie {
	cookie := createCookie(name, "")
	cookie.MaxAge = -1
	return cookie
}
