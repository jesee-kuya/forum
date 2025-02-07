package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type StoreSession struct {
	Token, Email string
	UserId       int
	ExpiryTime   time.Time
}

var (
	Session  StoreSession
	Sessions []StoreSession
)

func CookieSession(w http.ResponseWriter, session StoreSession) {
	paths := []string{"/home", "/upload", "/reaction", "/logout", "/comments", "/like", "/dislike"}

	for _, v := range paths {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    session.Token,
			Expires:  session.ExpiryTime,
			HttpOnly: true,
			Path:     v,
		})
	}
}

func ValidateCookie(r *http.Request) (StoreSession, *http.Cookie, error) {
	session := StoreSession{}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Cookie not found: %v", err)
		return StoreSession{}, nil, fmt.Errorf("cookie not found: %v", err)
	}

	for _, v := range Sessions {
		if v.Token == cookie.Value {
			session = v
			break
		}
	}
	return session, cookie, nil
}
