package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twsm000/lenslocked/controllers"
	"github.com/twsm000/lenslocked/models/database"
	"github.com/twsm000/lenslocked/models/database/postgres"
	"github.com/twsm000/lenslocked/models/repositories/postgresrepo"
	"github.com/twsm000/lenslocked/models/services"
	"github.com/twsm000/lenslocked/templates"
	"github.com/twsm000/lenslocked/views"

	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	logInfo  *log.Logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	logError *log.Logger = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
)

func main() {
	csrfAuthKey := flag.String("csrf-auth", "", "CSRF Auth Key (32 bytes - Mandatory)")
	secureCookie := flag.Bool("secure-cookie", true, "Secure the cookie when use https (CSRF Protection)")
	flag.Parse()
	csrfAuthKeyData := []byte(*csrfAuthKey)
	if len(csrfAuthKeyData) != 32 {
		log.Println("-csrf-auth needs to be 32 bytes")
		flag.PrintDefaults()
		os.Exit(1)
	}

	logInfo.Println("Secure-Cookie:", *secureCookie)
	db := MustGet(database.NewConnection(postgres.Config{
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

	server := http.Server{
		Addr:    ":8080",
		Handler: NewRouter(db, csrfAuthKeyData, *secureCookie),
	}
	Run(&server)
}

func ApplyHTML(page ...string) []string {
	return append([]string{"layout.tailwind.html", "footer.html"}, page...)
}

func NewRouter(db *sql.DB, csrfAuthKey []byte, secureCookie bool) http.Handler {
	homeTemplate := MustGet(views.ParseFSTemplate(templates.FS, ApplyHTML("home.html")...))
	contactTemplate := MustGet(views.ParseFSTemplate(templates.FS, ApplyHTML("contact.html")...))
	faqTemplate := MustGet(views.ParseFSTemplate(templates.FS, ApplyHTML("faq.html")...))
	signupTemplate := MustGet(views.ParseFSTemplate(templates.FS, ApplyHTML("signup.html")...))
	signinTemplate := MustGet(views.ParseFSTemplate(templates.FS, ApplyHTML("signin.html")...))

	userController := controllers.User{
		LogInfo:  logInfo,
		LogError: logError,
		UserService: services.User{
			Repository: MustGet(postgresrepo.NewUserRepository(db)),
		},
	}
	userController.Templates.SignUpPage = signupTemplate
	userController.Templates.SignInPage = signinTemplate
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", AsHTML(controllers.StaticTemplateHandler(homeTemplate)))
	router.Get("/contact", AsHTML(controllers.StaticTemplateHandler(contactTemplate)))
	router.Get("/faq", AsHTML(controllers.FAQ(faqTemplate)))
	router.Get("/signup", AsHTML(userController.SignUpPageHandler))
	router.Get("/signin", AsHTML(userController.SignInPageHandler))
	router.Post("/signin", AsHTML(userController.Authenticate))
	router.NotFound(AsHTML(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}))

	router.Route("/users", func(r chi.Router) {
		r.Post("/", userController.Create)
	})

	csrfMiddleware := csrf.Protect(
		csrfAuthKey,
		csrf.Secure(secureCookie),
	)

	return csrfMiddleware(router)
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

func MustGet[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
