package main

import (
	"context"
	"database/sql"
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
	"github.com/twsm000/lenslocked/models/entities"
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
	csrfAuthKey := flag.String("csrf-auth", "", "CSRF Auth Key (32 bytes - Mandatory)")
	secureCookie := flag.Bool("secure-cookie", true, "Secure the cookie when use https (CSRF Protection)")
	sessionTokenSize := flag.Int("session-token-size", entities.MinBytesPerSessionToken, "Size in bytes of the session tokens (Default 32 (min))")
	flag.Parse()
	csrfAuthKeyData := []byte(*csrfAuthKey)
	if len(csrfAuthKeyData) != 32 {
		log.Println("-csrf-auth needs to be 32 bytes")
		flag.PrintDefaults()
		os.Exit(1)
	}

	logInfo.Println("Secure-Cookie:", *secureCookie)
	logInfo.Println("SessionTokenSize", *sessionTokenSize)
	db := result.MustGet(database.NewConnection(postgres.Config{
		Driver:   "pgx",
		Host:     "localhost",
		Port:     5432,
		User:     "lenslocked",
		Password: "lenslocked",
		Database: "lenslocked",
		SSLMode:  "disable",
	}))
	defer func() {
		logInfo.Println("Closing database...")
		if err := db.Close(); err != nil {
			logError.Println("Database close error:", err)
		}
	}()
	TryTerminate(postgres.MigrateFS(db, "", migrations.FS))

	router, closer := NewRouter(db, csrfAuthKeyData, *secureCookie, *sessionTokenSize)
	defer func() {
		logInfo.Println("Closing resources...")
		if err := closer.Close(); err != nil {
			logError.Println(err)
		}
	}()
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	Run(&server)
}

func ApplyHTML(page ...string) []string {
	return append([]string{"layout.tailwind.html", "footer.html"}, page...)
}

func NewRouter(
	db *sql.DB,
	csrfAuthKey []byte,
	secureCookie bool,
	bytesPerToken int,
) (http.Handler, io.Closer) {
	homeTemplate := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("home.html")...))
	contactTemplate := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("contact.html")...))
	faqTemplate := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("faq.html")...))
	signupTemplate := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("signup.html")...))
	signinTemplate := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("signin.html")...))
	IntrnSrvErrTemplate := result.MustGet(views.ParseFSTemplate(logError, templates.FS, ApplyHTML("500.html")...))

	userRepo := result.MustGet(postgresrepo.NewUserRepository(db))
	sessionRepo := result.MustGet(postgresrepo.NewSessionRepository(db, logError, logInfo, logWarn))
	sessionService := services.NewSession(bytesPerToken, sessionRepo)
	userController := controllers.User{
		LogInfo:  logInfo,
		LogError: logError,
		UserService: services.User{
			Repository: userRepo,
		},
		SessionService: sessionService,
	}
	userController.Templates.SignUpPage = signupTemplate
	userController.Templates.SignInPage = signinTemplate

	csrfMiddleware := csrf.Protect(
		csrfAuthKey,
		csrf.Secure(secureCookie),
	)
	userMiddleware := controllers.UserMiddleware{
		LogWarn:        logWarn,
		SessionService: sessionService,
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(csrfMiddleware)
	router.Use(userMiddleware.SetUserToRequestContext)

	router.Get("/", AsHTML(controllers.StaticTemplateHandler(homeTemplate)))
	router.Get("/contact", AsHTML(controllers.StaticTemplateHandler(contactTemplate)))
	router.Get("/faq", AsHTML(controllers.FAQ(faqTemplate)))
	router.Get("/signup", AsHTML(userController.SignUpPageHandler))
	router.Get("/signin", AsHTML(userController.SignInPageHandler))
	router.Post("/signin", AsHTML(userController.Authenticate))
	router.Post("/signout", AsHTML(userController.SignOut))
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
		IntrnSrvErrTemplate.Execute(w, r, nil)
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
