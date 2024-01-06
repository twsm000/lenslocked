package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/twsm000/lenslocked/models/contextutil"
	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/httpll"
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
	SessionService services.Session
}

func (uc User) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	uc.Templates.SignUpPage.Execute(w, r, data)
}

func (uc User) SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	uc.Templates.SignInPage.Execute(w, r, data)
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
			httpll.SendStatusInternalServerError(w, r)

		default:
			err = errors.Join(ErrUntracked, err)
			httpll.SendStatusInternalServerError(w, r)
		}
		uc.LogError.Println(err)
		return
	}

	uc.LogInfo.Println("User created:", user)
	session, err := uc.SessionService.Create(user.ID)
	if err != nil {
		// TODO: validate other error types
		uc.LogError.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	uc.LogInfo.Println("Session created:", session)
	uc.createSessionCookieAndRedirect(w, r, session)
}

func (uc User) Authenticate(w http.ResponseWriter, r *http.Request) {
	var authCredentials entities.UserAuthenticable
	authCredentials.Email.Set(r.PostFormValue("email"))
	authCredentials.Password.Set(r.PostFormValue("password"))

	user, err := uc.UserService.Authenticate(authCredentials)
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrInvalidUserPassword),
			errors.Is(err, services.ErrInvalidAuthCredentials):
			url := fmt.Sprintf("/signin?email=%s", authCredentials.Email)
			http.Redirect(w, r, url, http.StatusFound)

		default:
			err = errors.Join(ErrUntracked, err)
			httpll.SendStatusInternalServerError(w, r)
		}

		uc.LogError.Println(err)
		return
	}

	uc.LogInfo.Println("User authenticated:", user)
	session, err := uc.SessionService.Create(user.ID)
	if err != nil {
		uc.LogError.Println(err)
		httpll.Redirect500Page(w, r)
		return
	}

	uc.LogInfo.Println("Session created:", session)
	uc.createSessionCookieAndRedirect(w, r, session)
}

func (uc *User) SignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = uc.SessionService.DeleteByToken(cookie.Value)
	if err != nil {
		uc.LogError.Println(err)
		httpll.Redirect500Page(w, r)
		return
	}

	http.SetCookie(w, deleteCookie(CookieSession))
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (uc *User) UserInfo(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUser(r.Context())
	if !ok {
		uc.LogInfo.Println("user not found in the current context")
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "User: %+v\n", user)
	fmt.Fprintf(w, "Header: %+v\n", r.Header)
}

func (uc *User) createSessionCookieAndRedirect(
	w http.ResponseWriter,
	r *http.Request,
	session *entities.Session,
) {
	cookie := createSessionCookie(session)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

type UserMiddleware struct {
	SessionService services.Session
}

func (um UserMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := um.SessionService.FindUserByToken(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(contextutil.WithUser(r.Context(), user)))
	})
}
