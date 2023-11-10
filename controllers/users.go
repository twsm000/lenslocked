package controllers

import (
	"net/http"

	"github.com/twsm000/lenslocked/views"
)

type Users struct {
	Templates struct {
		SignUpPage *views.Template
	}
}

func (u Users) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	u.Templates.SignUpPage.Execute(w, nil)
}
