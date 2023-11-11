package controllers

import (
	"fmt"
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

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	fmt.Fprintf(w, "Email: %s, Password: %s\n", email, password)
}
