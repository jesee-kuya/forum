package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	GoogleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL = "https://oauth2.googleapis.com/token"
	GoogleUserInfo = "https://www.googleapis.com/oauth2/v3/userinfo"
)

var (
	GoogleClientID     = os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
)

type GoogleUser struct {
	Sub, Name, Email string
}

// GoogleSignUp initiates the Google sign-up process by redirecting the user to Google's OAuth 2.0 server for authentication.
func GoogleSignUp(w http.ResponseWriter, r *http.Request) {
	// Generate a random state and set it as a cookie to prevent CSRF attacks
	state := generateStateCookie(w)

	// Construct the Google OAuth 2.0 authorization URL with necessary parameters
	redirectURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid email profile&state=%s&prompt=select_account&access_type=offline",
		GoogleAuthURL,
		GoogleClientID,
		url.QueryEscape("http://localhost:9000/auth/google/callback"),
		state,
	)

	// Redirect the user to Google's OAuth 2.0 server
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the callback from the Google OAuth 2.0 server after the user has granted the necessary permissions.
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Validate the state to prevent CSRF attacks
	if err := validateState(r); err != nil {
		log.Printf("State validation failed: %v", err)
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// Get the authorization code from the query parameter
	code := r.URL.Query().Get("code")
	token, err := exchangeGoogleToken(code)
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

	// Attempt to create or authenticate the user
	if handleUserAuth(w, user.Email, user.Name) {
		log.Printf("User authentication successful")
		http.Redirect(w, r, "/sign-in", http.StatusTemporaryRedirect)
	} else {
		log.Printf("User authentication failed")
		http.Redirect(w, r, "/sign-in?error=auth_failed", http.StatusTemporaryRedirect)
	}
}

// exchangeGoogleToken exchanges the authorization code for an access token from Google.
func exchangeGoogleToken(code string) (string, error) {
	// Prepare the data for the token exchange request
	data := url.Values{
		"code":          {code},
		"client_id":     {GoogleClientID},
		"client_secret": {GoogleClientSecret},
		"redirect_uri":  {"http://localhost:9000/auth/google/callback"},
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
