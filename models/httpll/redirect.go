package httpll

import (
	"fmt"
	"net/http"
)

func SendStatusInternalServerError(w http.ResponseWriter, r *http.Request) {
	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func Redirect500Page(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "redirect",
		Value:    fmt.Sprintf("%d", http.StatusInternalServerError),
		Path:     "/500",
		HttpOnly: true,
		MaxAge:   1,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/500", http.StatusSeeOther)
}
