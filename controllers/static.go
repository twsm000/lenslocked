package controllers

import (
	"html/template"
	"net/http"
)

func StaticTemplateHandler(tmpl Template[any]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, r, nil)
	}
}

func FAQ(tmpl Template[any]) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! We offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "What are your support hours?",
			Answer:   "We have support staff answering 24/7, though response times may be a bit slower on weekends",
		},
		{
			Question: "How I contact support?",
			Answer:   `Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>`,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, r, questions)
	}
}
