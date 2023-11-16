package controllers

import (
	"fmt"
	"net/http"

	"github.com/twsm000/lenslocked/models/services"
)

type User struct {
	Templates struct {
		SignUpPage Template
	}
	UserService services.User
}

func (u User) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignUpPage.Execute(w, data)
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	fmt.Fprintf(w, "Email: %s, Password: %s\n", email, password)
}
