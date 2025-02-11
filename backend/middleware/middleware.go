package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/repositories"
)

type contextSession string

var newSession contextSession = "userId"

// Authenticate middleware to check session token
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Println("NO session token", err)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		userID, err := repositories.ValidateSession(cookie.Value)
		if err != nil {
			log.Printf("Invalid session token: %v", err)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), newSession, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
