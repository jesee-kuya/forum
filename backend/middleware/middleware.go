package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/utils"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := repositories.ValidateSession(token)
			if err == nil {
				ctx := context.WithValue(r.Context(), "userID", userID)
				next(w, r.WithContext(ctx))
				return
			}
		}

		cookie, err := r.Cookie("session_token")
		if err == nil {
			userID, err := repositories.ValidateSession(cookie.Value)
			if err == nil {
				ctx := context.WithValue(r.Context(), "userID", userID)
				next(w, r.WithContext(ctx))
				return
			}
		}

		utils.ErrorHandler(w, "Unauthorized", http.StatusUnauthorized)
	}
}
