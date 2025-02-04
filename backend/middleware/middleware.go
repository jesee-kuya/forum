package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

// Authenticate middleware to check session token
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(repositories.Session.Token)
		if err != nil {
			util.ErrorHandler(w, "Unauthorized: No session token", http.StatusUnauthorized)
			return
		}

		userID, err := repositories.ValidateSession(cookie.Value)
		if err != nil {
			log.Printf("Invalid session token: %v", err)
			util.ErrorHandler(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		x := strconv.Itoa(userID)

		r.Header.Set("X-User-ID", x)

		next(w, r)
	}
}
