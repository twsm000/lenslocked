package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/repositories"
)

type User struct {
	LogInfo   *log.Logger
	LogError  *log.Logger
	Templates struct {
		SignUpPage Template
	}
	UserService interface {
		Create(input entities.UserCreatable) (*entities.User, error)
	}
}

func (uc User) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	uc.Templates.SignUpPage.Execute(w, data)
}

func (uc User) Create(w http.ResponseWriter, r *http.Request) {
	user, err := uc.UserService.Create(entities.UserCreatable{
		Email:    r.PostFormValue("email"),
		Password: []byte(r.PostFormValue("password")),
	})
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrInvalidUserEmail) ||
			errors.Is(err, entities.ErrInvalidUserPassword):
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		case errors.Is(err, entities.ErrInvalidUser) ||
			errors.Is(err, entities.ErrUserPasswordHash) ||
			errors.Is(err, repositories.ErrFailedToCreateUser):
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			err = errors.Join(errors.New("untracked error"), err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		uc.LogError.Println(err)
		return
	}

	fmt.Fprintf(w, "user created: %+v\n", user)
}
