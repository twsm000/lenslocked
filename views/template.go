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
	"sync"

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
	tmpl = tmpl.Funcs(template.FuncMap{
		"CSRFField": func(req *http.Request) (template.HTML, error) {
			return "", ErrNotImplemented
		},
	})
	tmpl, err := tmpl.ParseFS(fs, pattern...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fs template: %w", err)
	}
	return &Template[T]{
		logError: logError,
		htmlTmpl: tmpl,
	}, nil
}

type Template[T any] struct {
	syncOnceFuncs sync.Once
	htmlTmpl      *template.Template
	logError      *log.Logger
}

func (t *Template[T]) Execute(w http.ResponseWriter, r *http.Request, data T, errors ...entities.ClientError) {
	t.syncOnceFuncs.Do(func() {
		t.htmlTmpl = t.htmlTmpl.Funcs(template.FuncMap{
			"CSRFField": func(req *http.Request) template.HTML {
				return csrf.TemplateField(req)
			},
		})
	})
	var reqData struct {
		HTTPRequest *http.Request
		User        *entities.User
		Data        T
		Errors      []string
	}
	reqData.Data = data
	reqData.User = result.ExtractValue(contextutil.GetUser(r.Context()))
	reqData.HTTPRequest = r
	if len(errors) > 0 {
		reqData.Errors = make([]string, 0, len(errors))
		for _, err := range errors {
			reqData.Errors = append(reqData.Errors, err.ClientErr())
		}
	}
	var buf bytes.Buffer
	if err := t.htmlTmpl.Execute(&buf, reqData); err != nil {
		t.logError.Println("Failed to execute template:", err)
		httpll.Redirect500Page(w, r)
		return
	}

	_, err := io.Copy(w, &buf)
	if err != nil {
		t.logError.Println("Failed to send data to ResponseWriter:", err)
	}
}
