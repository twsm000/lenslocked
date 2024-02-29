package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/twsm000/lenslocked/models/contextutil"
	"github.com/twsm000/lenslocked/models/entities"
	"github.com/twsm000/lenslocked/models/httpll"
	"github.com/twsm000/lenslocked/pkg/result"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

func ParseFSTemplate[T any](logError *log.Logger, fs fs.FS, pattern ...string) (*Template[T], error) {
	tmpl := template.New(pattern[0])
	tmpl, err := tmpl.ParseFS(fs, pattern...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fs template: %w", err)
	}
	return &Template[T]{
		logError: logError,
		htmlTmpl: tmpl,
	}, nil
}

type templateData[T any] struct {
	CSRFField template.HTML
	User      *entities.User
	Data      T
	Errors    []string
}
type Template[T any] struct {
	htmlTmpl *template.Template
	logError *log.Logger
}

func (t *Template[T]) Execute(w http.ResponseWriter, r *http.Request, data T, errors ...entities.ClientError) {
	tmplData := templateData[T]{
		CSRFField: csrf.TemplateField(r),
		Data:      data,
		User:      result.ExtractValue(contextutil.GetUser(r.Context())),
		Errors:    toStringSlice(errors),
	}

	var buf bytes.Buffer
	if err := t.htmlTmpl.Execute(&buf, tmplData); err != nil {
		t.logError.Println("Failed to execute template:", err)
		httpll.Redirect500Page(w, r)
		return
	}

	_, err := io.Copy(w, &buf)
	if err != nil {
		t.logError.Println("Failed to send data to ResponseWriter:", err)
	}
}

func toStringSlice(errors []entities.ClientError) []string {
	var result []string
	if len(errors) > 0 {
		result = make([]string, 0, len(errors))
		for _, err := range errors {
			result = append(result, err.ClientErr())
		}
	}
	return result
}
