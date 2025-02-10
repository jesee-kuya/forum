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

	cookie, err := getSessionID(r)
	if err != nil {
		log.Println("Invalid Session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = repositories.DeleteSession(cookie)
	if err != nil {
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}
	delete(SessionStore, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
