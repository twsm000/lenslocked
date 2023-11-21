package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
	"github.com/twsm000/lenslocked/models/services"
)

type User struct {
	LogInfo   *log.Logger
	LogError  *log.Logger
	Templates struct {
		SignUpPage Template
		SignInPage Template
	}
	UserService interface {
		Create(input entities.UserCreatable) (*entities.User, error)
		Authenticate(input entities.UserAuthenticable) (*entities.User, error)
	}
}

func (uc User) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email     string
		CSRFField template.HTML
	}
	data.Email = r.FormValue("email")
	data.CSRFField = csrf.TemplateField(r)
	uc.Templates.SignUpPage.Execute(w, data)
}

func (uc User) SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	uc.Templates.SignInPage.Execute(w, data)
}

func (uc User) Create(w http.ResponseWriter, r *http.Request) {
	var userInput entities.UserCreatable
	userInput.Email.Set(r.PostFormValue("email"))
	userInput.Password.Set(r.PostFormValue("password"))

	user, err := uc.UserService.Create(userInput)
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrInvalidUserEmail) ||
			errors.Is(err, entities.ErrInvalidUserPassword):
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		case errors.Is(err, entities.ErrInvalidUser) ||
			errors.Is(err, entities.ErrFailedToHashPassword) ||
			errors.Is(err, repositories.ErrFailedToCreateUser):
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		default:
			err = errors.Join(ErrUntracked, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		uc.LogError.Println(err)
		return
	}

	fmt.Fprintf(w, "user created: %+v\n", user)
}

func (uc User) Authenticate(w http.ResponseWriter, r *http.Request) {
	var authCredentials entities.UserAuthenticable
	authCredentials.Email.Set(r.PostFormValue("email"))
	authCredentials.Password.Set(r.PostFormValue("password"))

	user, err := uc.UserService.Authenticate(authCredentials)
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrInvalidUserPassword):
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		case errors.Is(err, services.ErrInvalidAuthCredentials):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		default:
			err = errors.Join(ErrUntracked, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		uc.LogError.Println(err)
		return
	}

	cookie := http.Cookie{
		Name:     "email",
		Value:    user.Email.String(),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "%+v", user)
}
