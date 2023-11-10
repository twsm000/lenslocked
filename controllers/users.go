package controllers

import (
	"net/http"
)

type Users struct {
	Templates struct {
		SignUpPage Template
	}
}

func (u Users) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	u.Templates.SignUpPage.Execute(w, nil)
}
