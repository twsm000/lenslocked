package entities

import (
	"fmt"
	"strings"
)

type Email struct {
	email string
}

func (e *Email) Set(email string) {
	e.email = strings.ToLower(email)
}

func (e Email) IsEmpty() bool {
	return e.email == ""
}

func (e Email) String() string {
	return e.email
}

func (e *Email) Scan(value any) error {
	if value == nil {
		e.email = ""
		return nil
	}

	email, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid scan type: %T", value)
	}
	e.Set(email)
	return nil
}
