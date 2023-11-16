package main

import (
	"context"
	"database/sql"
	"errors"
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
	"github.com/twsm000/lenslocked/templates"
	"github.com/twsm000/lenslocked/views"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	logInfo  *log.Logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	logError *log.Logger = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
)

func main() {
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
		Handler: NewRouter(db),
	}
	Run(&server)
}

func ApplyHTML(page ...string) []string {
	return append([]string{"layout.tailwind.html", "footer.html"}, page...)
}

func NewRouter(db *sql.DB) http.Handler {
	homeTemplate := MustGet(views.ParseFSTemplate(templates.FS, ApplyHTML("home.html")...))
	contactTemplate := MustGet(views.ParseFSTemplate(templates.FS, ApplyHTML("contact.html")...))
	faqTemplate := MustGet(views.ParseFSTemplate(templates.FS, ApplyHTML("faq.html")...))
	signupTemplate := MustGet(views.ParseFSTemplate(templates.FS, ApplyHTML("signup.html")...))

	usersController := controllers.Users{}
	usersController.Templates.SignUpPage = signupTemplate
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", AsHTML(controllers.StaticTemplateHandler(homeTemplate)))
	router.Get("/contact", AsHTML(controllers.StaticTemplateHandler(contactTemplate)))
	router.Get("/faq", AsHTML(controllers.FAQ(faqTemplate)))
	router.Get("/signup", AsHTML(usersController.SignUpPageHandler))
	router.NotFound(AsHTML(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}))

	router.Route("/users", func(r chi.Router) {
		r.Post("/", usersController.Create)
	})

	return router
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

func executeTemplate(w http.ResponseWriter, fpath string) {
	tmpl, err := views.ParseTemplate(fpath)
	if err != nil {
		logError.Println("failed to parse error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
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
