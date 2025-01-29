package middleware

import (
	"context"
	"net/http"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/utils"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.ErrorHandler(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			utils.ErrorHandler(w, "Bad Request", http.StatusBadRequest)
			return
		}

		sessionToken := cookie.Value
		userID, err := repositories.ValidateSession(sessionToken)
		if err != nil {
			utils.ErrorHandler(w, "Invalid session or expired token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
