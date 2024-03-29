package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/twsm000/lenslocked/models/contextutil"
	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/httpll"
	"github.com/twsm000/lenslocked/models/repositories"
	"github.com/twsm000/lenslocked/models/services"
)

type SignUpPageData struct {
	Email string
}

type SignInPageData struct {
	Email string
}

type User struct {
	LogInfo   *log.Logger
	LogError  *log.Logger
	Templates struct {
		SignUpPage            Template[SignUpPageData]
		SignInPage            Template[SignInPageData]
		ForgotPasswordPage    Template[any]
		CheckPasswordSentPage Template[any]
		ResetPasswordPage     Template[any]
	}
	UserService          services.User
	SessionService       services.Session
	PasswordResetService services.PasswordReset
	EmailService         *services.EmailService
}

func (uc *User) SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	uc.Templates.SignUpPage.Execute(w, r, SignUpPageData{r.FormValue("email")})
}

func (uc *User) SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	uc.Templates.SignInPage.Execute(w, r, SignInPageData{r.FormValue("email")})
}

func (uc *User) ForgotPasswordPageHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	uc.Templates.ForgotPasswordPage.Execute(w, r, data)
}

func (uc *User) ResetPasswordPageHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	uc.Templates.ResetPasswordPage.Execute(w, r, data)
}

func (uc *User) Create(w http.ResponseWriter, r *http.Request) {
	var userInput entities.UserCreatable
	userInput.Email.Set(r.PostFormValue("email"))
	userInput.Password.Set(r.PostFormValue("password"))
	signUpPageData := SignUpPageData{Email: r.PostFormValue("email")}

	user, err := uc.UserService.Create(userInput)
	if err != nil {
		uc.LogError.Println(err)
		if !err.IsClientErr() {
			err = entities.NewClientError("Sorry, something went wrong when creating your user account", err)
		}
		uc.Templates.SignUpPage.Execute(w, r, signUpPageData, err)
		return
	}

	uc.LogInfo.Println("User created:", user)
	session, err := uc.SessionService.Create(user.ID)
	if err != nil {
		uc.LogError.Println(err)
		if err.IsClientErr() {
			if err.Is(repositories.ErrUserNotFound) {
				err = entities.NewClientError("Something went wrong while creating your user, please try again!", err)
			}
			uc.Templates.SignUpPage.Execute(w, r, signUpPageData, err)
			return
		}

		httpll.Redirect500Page(w, r)
		return
	}

	uc.LogInfo.Println("Session created:", session)
	uc.createSessionCookieAndRedirect(w, r, session)
}

func (uc *User) Authenticate(w http.ResponseWriter, r *http.Request) {
	var authCredentials entities.UserAuthenticable
	authCredentials.Email.Set(r.PostFormValue("email"))
	authCredentials.Password.Set(r.PostFormValue("password"))
	signInPageData := SignInPageData{Email: r.PostFormValue("email")}

	user, err := uc.UserService.Authenticate(authCredentials)
	if err != nil {
		uc.LogError.Println(err)
		if err.IsClientErr() {
			uc.Templates.SignInPage.Execute(w, r, signInPageData, err)
			return
		}

		httpll.Redirect500Page(w, r)
		return
	}

	uc.LogInfo.Println("User authenticated:", user)
	session, err := uc.SessionService.Create(user.ID)
	if err != nil {
		uc.LogError.Println(err)
		if err.IsClientErr() {
			if err.Is(repositories.ErrUserNotFound) {
				err = entities.NewClientError("Something went wrong while authenticating your user, please try again!", err)
			}
			uc.Templates.SignInPage.Execute(w, r, signInPageData, err)
			return
		}

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

func (uc *User) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var email entities.Email
	email.Set(r.PostFormValue("email"))
	pr, err := uc.PasswordResetService.Create(email)
	if err != nil {
		// TODO: handle all the cases
		uc.LogError.Println(err)
		httpll.SendStatusInternalServerError(w, r)
		return
	}

	// TODO: generate reset url from the correct domain
	url := url.Values{
		"token": {
			pr.Token.Value(),
		},
	}
	resetURL := fmt.Sprintf("http:localhost:8080/resetpass?%s", url.Encode())
	err = uc.EmailService.ForgotPassword(email.String(), resetURL)
	if err != nil {
		// TODO: handle all the cases
		uc.LogError.Println(err)
		httpll.SendStatusInternalServerError(w, r)
		return
	}

	uc.Templates.CheckPasswordSentPage.Execute(w, r, struct{ Email string }{Email: email.String()})
}

func (uc *User) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var stoken entities.SessionToken
	stoken.SetFromHex(r.PostFormValue("token"))
	user, err := uc.PasswordResetService.Consume(stoken)
	if err != nil {
		// TODO: handle all the cases
		uc.LogError.Println(err)
		httpll.SendStatusInternalServerError(w, r)
		return
	}

	var rawPassword entities.RawPassword
	rawPassword.Set(r.PostFormValue("password"))
	if err := uc.UserService.UpdatePassword(user, rawPassword); err != nil {
		// TODO: handle all the cases
		uc.LogError.Println(err)
		httpll.SendStatusInternalServerError(w, r)
		return
	}

	uc.LogInfo.Println("User password updated:", user)
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

func (uc *User) UserInfo(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUser(r.Context())
	if !ok {
		uc.LogError.Println("Required user was not found in the current context.")
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "User: %+v\n", user)
	fmt.Fprintf(w, "Header: %+v\n", r.Header)
}

func (uc *User) createSessionCookieAndRedirect(w http.ResponseWriter, r *http.Request, session *entities.Session) {
	cookie := createSessionCookie(session)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

type UserMiddleware struct {
	LogWarn        *log.Logger
	SessionService services.Session
}

func (um UserMiddleware) SetUserToRequestContext(next http.Handler) http.Handler {
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

func (um UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := contextutil.GetUser(r.Context()); !ok {
			um.LogWarn.Println("user not found in the current request context")
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
