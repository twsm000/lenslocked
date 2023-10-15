package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", hello)

	log.Println("Starting server at port:", 8080)
	http.ListenAndServe(":8080", nil)
}
func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "Hello World! %s\n", time.Now())
}
