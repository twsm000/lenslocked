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

func ParseFSTemplate(logError *log.Logger, fs fs.FS, pattern ...string) (*Template, error) {
	tmpl := template.New(pattern[0])
	tmpl = tmpl.Funcs(template.FuncMap{
		"CSRFField": func(req *http.Request) (template.HTML, error) {
			return "", ErrNotImplemented
		},
		"errors": func() []string {
			return []string{
				"Don't do that!",
				"The email address you provided is already associated with an account.",
				"Something went wrong.",
			}
		},
	})
	tmpl, err := tmpl.ParseFS(fs, pattern...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fs template: %w", err)
	}
	return &Template{
		logError: logError,
		htmlTmpl: tmpl,
	}, nil
}

type Template struct {
	syncOnceFuncs sync.Once
	htmlTmpl      *template.Template
	logError      *log.Logger
}

func (t *Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
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
		Data        any
	}
	reqData.Data = data
	reqData.User = result.ExtractValue(contextutil.GetUser(r.Context()))
	reqData.HTTPRequest = r
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
