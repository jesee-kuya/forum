package middleware

import (
	"context"
	"net/http"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

// validates the session token and sets the user in the context
func SessionMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getSessionID(r)
		if err != nil {
			util.ErrorHandler(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := repositories.ValidateSession(sessionID)
		if err != nil {
			util.ErrorHandler(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
