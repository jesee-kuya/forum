package middleware

import (
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/repositories"
)

// Authenticate middleware to check session token
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Println("NO session token", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		userID, err := repositories.ValidateSession(cookie.Value)
		if err != nil {
			log.Printf("Invalid session token: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		r.Header.Set("X-User-ID", userID)

		next(w, r)
	}
}
