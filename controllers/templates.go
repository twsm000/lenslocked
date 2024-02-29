package controllers

import (
	"net/http"

	"github.com/twsm000/lenslocked/models/entities"
)

type Template[T any] interface {
	Execute(w http.ResponseWriter, r *http.Request, data T, errors ...entities.ClientError)
}
