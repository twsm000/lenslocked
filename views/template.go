package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

func ParseTemplate(fpath string) (*Template, error) {
	tmpl, err := template.ParseFiles(fpath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	return &Template{htmlTmpl: tmpl}, nil
}

func ParseFSTemplate(fs fs.FS, pattern ...string) (*Template, error) {
	tmpl, err := template.ParseFS(fs, pattern...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fs template: %w", err)
	}
	return &Template{htmlTmpl: tmpl}, nil
}

type Template struct {
	htmlTmpl *template.Template
}

func (t *Template) Execute(w http.ResponseWriter, data any) {
	if err := t.htmlTmpl.Execute(w, data); err != nil {
		log.Println("failed to execute template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
