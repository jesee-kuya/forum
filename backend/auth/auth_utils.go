package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/util"
)

/*
handleUserAuth attempts to authenticate a user by email. If the user is not found, handleUserAuth will create a new user and log them in. Returns true if the user is authenticated, false otherwise.
*/
func handleUserAuth(w http.ResponseWriter, email, username string) bool {
	var userID int
	err := util.DB.QueryRow(
		"SELECT id FROM tblUsers WHERE email = ?", email,
	).Scan(&userID)

	// Create a new user if not found in db
	if errors.Is(err, sql.ErrNoRows) {
		res, err := util.DB.Exec(
			"INSERT INTO tblUsers(username, email) VALUES(?, ?)",
			username, email,
		)
		if err != nil {
			log.Printf("User creation failed: %v", err)
			return false
		}
		id, _ := res.LastInsertId()
		userID = int(id)
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return false
	}

	// Set the session cookie for the user
	SetSessionCookie(w, userID)
	return true
}

// // SetSessionCookie sets a session cookie for the given user ID.
func SetSessionCookie(w http.ResponseWriter, userID int) {
	token := generateSessionToken()
	_, err := util.DB.Exec(
		"INSERT INTO tblSessions(user_id, session_token) VALUES(?, ?)",
		userID, token,
	)
	if err != nil {
		log.Printf("Session creation failed: %v", err)
		return
	}

	// Set the session cookie for the user
	http.SetCookie(w, &http.Cookie{
		Name:     "forum_session",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   86400, // 1 day
	})
}

/*
generateSessionToken generates a secure random session token. It returns the token encoded in URL-safe base64 format.
*/
func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)

	// Encode the random bytes into a URL-safe base64 string and return it
	return base64.URLEncoding.EncodeToString(b)
}
