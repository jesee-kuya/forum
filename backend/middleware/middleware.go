package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jesee-kuya/forum/backend/repositories"
)

// Authenticate middleware to check session token
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized: No session token", http.StatusUnauthorized)
			return
		}

		userID, err := repositories.ValidateSession(cookie.Value)
		if err != nil {
			log.Printf("Invalid session token: %v", err)
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		x := strconv.Itoa(userID)

		r.Header.Set("X-User-ID", x)

		next(w, r)
	}
}
