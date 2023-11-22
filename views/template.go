package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

func ParseFSTemplate(fs fs.FS, pattern ...string) (*Template, error) {
	tmpl := template.New(pattern[0])
	tmpl = tmpl.Funcs(template.FuncMap{
		"CSRFField": func() template.HTML {
			return `<input type="hidden"/>`
		},
	})
	tmpl, err := tmpl.ParseFS(fs, pattern...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fs template: %w", err)
	}
	return &Template{htmlTmpl: tmpl}, nil
}

type Template struct {
	htmlTmpl *template.Template
}

func (t *Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	if err := t.htmlTmpl.Execute(w, data); err != nil {
		log.Println("failed to execute template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
