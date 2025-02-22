package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/jesee-kuya/forum/backend/util"
)

// GoogleSignIn initiates the Google sign-in process by redirecting the user to Google's OAuth 2.0 server for authentication.
func GoogleSignIn(w http.ResponseWriter, r *http.Request) {
	// Generate a random state and set it as a cookie to prevent CSRF attacks
	state := generateStateCookie(w)

	// Construct the Google OAuth 2.0 authorization URL with necessary parameters
	redirectURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid email profile&state=%s&prompt=select_account&access_type=offline",
		GoogleAuthURL,
		util.GoogleClientID,
		url.QueryEscape("http://localhost:9000/auth/google/signin/callback"),
		state,
	)

	// Set the CORS header to allow the request to be made from the frontend
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:9000")

	// Redirect the user to Google's OAuth 2.0 server
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// GoogleSignInCallback handles the callback from the Google OAuth 2.0 server after the user has granted the necessary permissions.
func GoogleSignInCallback(w http.ResponseWriter, r *http.Request) {
	// Validate the state to prevent CSRF attacks
	if err := validateState(r); err != nil {
		log.Println("Invalid state")
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// Get the authorization code from the query parameter
	code := r.URL.Query().Get("code")

	// Exchange the authorization code for an access token
	token, err := exchangeGoogleTokenSignIn(code)
	if err != nil {
		log.Printf("Token exchange failed: %v\n", err)
		http.Redirect(w, r, "/sign-in?error=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	// Get the user information from the Google UserInfo endpoint
	user, err := getGoogleUser(token)
	if err != nil {
		log.Printf("Failed to get user info: %v\n", err)
		http.Redirect(w, r, "/sign-in?error=user_info_failed", http.StatusTemporaryRedirect)
		return
	}

	// Check if the user exists in the database
	var userID int
	err = util.DB.QueryRow("SELECT id FROM tblUsers WHERE email = ?", user.Email).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, "/sign-up?error=no_account", http.StatusTemporaryRedirect)
			return
		}
		log.Printf("Database error: %v", err)
		http.Redirect(w, r, "/sign-in?error=database_error", http.StatusTemporaryRedirect)
		return
	}

	setSessionCookie(w, userID)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

/*
exchangeGoogleTokenSignIn exchanges the authorization code for an access token from Google. It returns the access token that can be used to access the user's information.
*/
func exchangeGoogleTokenSignIn(code string) (string, error) {
	data := url.Values{
		"code":          {code},
		"client_id":     {util.GoogleClientID},
		"client_secret": {util.GoogleClientSecret},
		"redirect_uri":  {"http://localhost:9000/auth/google/signin/callback"},
		"grant_type":    {"authorization_code"},
	}

	// Make a POST request to Google's token URL
	resp, err := http.PostForm(GoogleTokenURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response to extract the access token
	var result struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.AccessToken, nil
}
