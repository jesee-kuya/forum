package handler

import (
	"net/http"
	"time"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorHandler(w, "No active session", http.StatusUnauthorized)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		util.ErrorHandler(w, "No active session", http.StatusUnauthorized)
		return
	}

	err = repositories.DeleteSession(cookie.Value)
	if err != nil {
		util.ErrorHandler(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	for i, session := range Sessions {
		if session.Token == cookie.Value {
			Sessions[i] = StoreSession{}
			break
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
