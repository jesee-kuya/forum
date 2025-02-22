package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/util"
)

// generateStateCookie generates a random state and sets it as a cookie.
func generateStateCookie(w http.ResponseWriter) string {
	// Generate a random string
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Printf("error generating random state: %v", err)
		return ""
	}

	state := base64.URLEncoding.EncodeToString(b)

	// Set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode,
	})
	return state
}

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
	setSessionCookie(w, userID)
	return true
}

// setSessionCookie sets a session cookie for the given user ID.
func setSessionCookie(w http.ResponseWriter, userID int) {
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
validateState checks if the state parameter in the URL matches the state parameter in the cookie. If the two do not match, validateState returns an error.
*/
func validateState(r *http.Request) error {
	state := r.URL.Query().Get("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		log.Printf("Cookie error: %v", err)
		return err
	}

	// Check if the two states match
	if cookie.Value != state {
		log.Printf("State mismatch. Cookie: %s, State: %s", cookie.Value, state)
		return errors.New("invalid state")
	}
	return nil
}

/*
getGoogleUser takes a token and makes a request to the Google UserInfo API to retrieve the user's information. getGoogleUser returns a GoogleUser struct if the request is successful, or an error if the request fails.
*/
func getGoogleUser(token string) (*GoogleUser, error) {
	req, _ := http.NewRequest("GET", GoogleUserInfo, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode the response as a GoogleUser
	var user GoogleUser
	json.NewDecoder(resp.Body).Decode(&user)

	// Return the user if the request was successful
	return &user, nil
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
