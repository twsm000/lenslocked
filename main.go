package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Printf("Starting server at port: http://localhost:%d\n", 8080)
	http.ListenAndServe(":8080", http.HandlerFunc(pathHandler))
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<h1>Home</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:twsm000@gmail.com\">twsm000@gmail.com</a>.</p>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
