package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	logInfo  *log.Logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	logError *log.Logger = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: NewRouter(),
	}
	Run(&server)
}

func NewRouter() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", setContentTypeTextHtml(homeHandler))
	router.Get("/contact", setContentTypeTextHtml(contactHandler))
	router.Get("/faq", setContentTypeTextHtml(faqHandler))
	router.NotFound(setContentTypeTextHtml(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}))

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
	tmpl, err := template.ParseFiles(fpath)
	if err != nil {
		log.Println("failed to parse error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("failed to execute template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "home.html"))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "contact.html"))
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<h1>FAQ Page</h1>
  <ul>
	<li>
	  <b>Is there a free version?</b>
	  Yes! We offer a free trial for 30 days on any paid plans.
	</li>
	<li>
	  <b>What are your support hours?</b>
	  We have support staff answering emails 24/7, though response
	  times may be a bit slower on weekends.
	</li>
	<li>
	  <b>How do I contact support?</b>
	  Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>
	</li>
  </ul>
  `)
}

func setContentTypeTextHtml(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next(w, r)
	}
}
