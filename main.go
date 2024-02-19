package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twsm000/lenslocked/controllers"
	"github.com/twsm000/lenslocked/models/database"
	"github.com/twsm000/lenslocked/models/database/postgres"
	"github.com/twsm000/lenslocked/models/repositories/postgresrepo"
	"github.com/twsm000/lenslocked/models/services"
	"github.com/twsm000/lenslocked/models/sql/postgres/migrations"
	"github.com/twsm000/lenslocked/pkg/result"
	"github.com/twsm000/lenslocked/templates"
	"github.com/twsm000/lenslocked/views"

	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	logError *log.Logger = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Llongfile)
	logInfo  *log.Logger = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Llongfile)
	logWarn  *log.Logger = log.New(os.Stdout, "WARN: ", log.LstdFlags|log.Llongfile)
)

func main() {
	envFilePath := flag.String("env-file", "", "Environment file settings")
	flag.Parse()

	env := result.MustGet(LoadEnvSettings(*envFilePath, "pgx"))
	if len(env.CSRF.Key) != 32 {
		log.Println("CSRF.Key needs to be 32 bytes")
		flag.PrintDefaults()
		os.Exit(1)
	}
	logInfo.Printf("EnvSettings: %s\n", result.MustGet(json.MarshalIndent(&env, "", "  ")))

	db := result.MustGet(database.NewConnection(env.DBConfig))
	defer func() {
		logInfo.Println("Closing database...")
		if err := db.Close(); err != nil {
			logError.Println("Database close error:", err)
		}
	}()
	TryTerminate(postgres.MigrateFS(db, "", migrations.FS))

	router, closer := NewRouter(db, env)
	defer func() {
		logInfo.Println("Closing resources...")
		if err := closer.Close(); err != nil {
			logError.Println(err)
		}
	}()
	server := http.Server{
		Addr:    env.Server.Address,
		Handler: router,
	}
	Run(&server)
}

func ApplyHTML(page ...string) []string {
	return append([]string{"layout.tailwind.html", "footer.html"}, page...)
}

func NewRouter(DB *sql.DB, env *EnvConfig) (http.Handler, io.Closer) {
	homeTmpl := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("home.html")...))
	contactTmpl := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("contact.html")...))
	faqTmpl := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("faq.html")...))
	signupTmpl := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("signup.html")...))
	signinTmpl := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("signin.html")...))
	forgotPasswordTmpl := result.MustGet(views.ParseFSTemplate(logError, templates.FS,
		ApplyHTML("forgot_password.html")...))
	checkPasswordSentTmpl := result.MustGet(views.ParseFSTemplate(logError, templates.FS,
		ApplyHTML("check_password_sent.html")...))
	IntrnSrvErrTmpl := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("500.html")...))

	userRepo := result.MustGet(postgresrepo.NewUserRepository(DB))
	userService := services.User{Repository: userRepo}
	sessionRepo := result.MustGet(postgresrepo.NewSessionRepository(DB, logError, logInfo, logWarn))
	sessionService := services.NewSession(env.Session.TokenSize, sessionRepo)
	passwordResetRepo := result.MustGet(postgresrepo.NewPasswordResetRepository(DB, logError, logInfo, logWarn))
	passwordResetService := services.NewPasswordReset(
		env.Session.TokenSize,
		services.DefaultPasswordResetDuration, // TODO: load this value from env file
		passwordResetRepo,
		userService,
	)
	emailService := services.NewEmailService(env.SMTPConfig)

	userController := controllers.User{
		LogInfo:              logInfo,
		LogError:             logError,
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}
	userController.Templates.SignUpPage = signupTmpl
	userController.Templates.SignInPage = signinTmpl
	userController.Templates.ForgotPasswordPage = forgotPasswordTmpl
	userController.Templates.CheckPasswordSentPage = checkPasswordSentTmpl

	csrfMiddleware := csrf.Protect([]byte(env.CSRF.Key), csrf.Secure(env.CSRF.Secure))
	userMiddleware := controllers.UserMiddleware{
		LogWarn:        logWarn,
		SessionService: sessionService,
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(csrfMiddleware)
	router.Use(userMiddleware.SetUserToRequestContext)

	router.Get("/", AsHTML(controllers.StaticTemplateHandler(homeTmpl)))
	router.Get("/contact", AsHTML(controllers.StaticTemplateHandler(contactTmpl)))
	router.Get("/faq", AsHTML(controllers.FAQ(faqTmpl)))
	router.Get("/signup", AsHTML(userController.SignUpPageHandler))
	router.Get("/signin", AsHTML(userController.SignInPageHandler))
	router.Get("/forgotpass", AsHTML(userController.ForgotPasswordHandler))
	router.Post("/signin", AsHTML(userController.Authenticate))
	router.Post("/signout", AsHTML(userController.SignOut))
	router.Post("/resetpass", AsHTML(userController.ResetPassword))
	router.NotFound(AsHTML(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}))
	router.Get("/500", AsHTML(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("redirect")
		if err != nil ||
			result.ExtractValue(strconv.Atoi(cookie.Value)) != http.StatusInternalServerError {
			if err == nil {
				logInfo.Printf("Cookie: %+v", cookie)
			}
			logError.Println("GET /500 EXPECT NOT FOUND")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		IntrnSrvErrTmpl.Execute(w, r, nil)
	}))

	router.Route("/users", func(r chi.Router) {
		r.Post("/", userController.Create)

		r.Route("/me", func(r chi.Router) {
			r.Use(userMiddleware.RequireUser)
			r.Get("/", userController.UserInfo)
		})
	})

	closer := func() error {
		return errors.Join(
			userRepo.Close(),
			sessionRepo.Close(),
		)
	}

	return router, CloserFunc(closer)
}

func Run(server *http.Server) {
	go func() {
		logInfo.Printf("Starting server at port: %s\n", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logError.Println("Failed to close http server:", err)
		}
	}()

	gracefullyShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefullyShutdown, syscall.SIGINT, syscall.SIGTERM)
	<-gracefullyShutdown
	logInfo.Println("Closing http server gracefully.")
	shutdownServerCtx, closeServer := context.WithTimeout(context.Background(), 10*time.Second)
	defer closeServer()
	if err := server.Shutdown(shutdownServerCtx); err != nil {
		logError.Println("Failed to shutdown gracefully http server:", err)
	}
	fmt.Println("Bye...")
}

// AsHTML set the response header Content-Type to text/html
func AsHTML(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next(w, r)
	}
}

type CloserFunc func() error

func (cf CloserFunc) Close() error {
	return cf()
}

func TryTerminate(err error) {
	if err != nil {
		logError.Fatalln(err)
	}
}
