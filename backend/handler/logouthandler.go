package handler

import (
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorHandler(w, "No active session", http.StatusUnauthorized)
		return
	}

	_, cookie, err := ValidateCookie(r)
	if err != nil {
		log.Printf("Failed to validate cookie: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
