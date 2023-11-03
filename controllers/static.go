package controllers

import (
	"net/http"

	"github.com/twsm000/lenslocked/views"
)

func StaticTemplateHandler(tmpl *views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}
