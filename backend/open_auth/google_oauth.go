package openauth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jesee-kuya/forum/backend/handler"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

const (
	GoogleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL = "https://oauth2.googleapis.com/token"
	GoogleUserInfo = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type GoogleUser struct {
	Sub, Name, Email string
}

// GoogleAuth initiates the Google authentication process (signup or signin)
func GoogleAuth(w http.ResponseWriter, r *http.Request) {
	state := generateStateCookie(w)

	// Construct the Google OAuth 2.0 authorization URL with necessary parameters
	redirectURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid email profile&state=%s&prompt=select_account&access_type=offline",
		GoogleAuthURL,
		util.GoogleClientID,
		url.QueryEscape(RedirectBaseURL+"/auth/google/callback"),
		state,
	)

	w.Header().Set("Access-Control-Allow-Origin", RedirectBaseURL)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the callback from Google's OAuth server
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
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

	user, err := getGoogleUser(token)
	if err != nil {
		log.Printf("Failed to get user info: %v\n", err)
		http.Redirect(w, r, "/sign-in?error=user_info_failed", http.StatusTemporaryRedirect)
		return
	}

	var (
		userID       int
		authProvider string
	)
	err = util.DB.QueryRow("SELECT id, auth_provider FROM tblUsers WHERE email = ?", user.Email).Scan(&userID, &authProvider)

	// If the email exists but with a different provider
	if err == nil && authProvider != "google" {
		log.Printf("Email already registered with %s: %v", authProvider, user.Email)
		http.Redirect(w, r, "/sign-in?error=email_exists&provider="+authProvider, http.StatusTemporaryRedirect)
		return
	}

	isNewUser := false

	// If user doesn't exist, create a new one
	if errors.Is(err, sql.ErrNoRows) {
		var count int
		err = util.DB.QueryRow("SELECT COUNT(*) FROM tblUsers WHERE username = ?", user.Name).Scan(&count)
		if err != nil {
			log.Printf("Database error checking username: %v", err)
			http.Redirect(w, r, "/sign-in?error=database_error", http.StatusTemporaryRedirect)
			return
		}

		if count > 0 {
			// Username is taken, generate a unique one by appending a random suffix
			user.Name = fmt.Sprintf("%s_%s", user.Name, user.Sub[:6])
		}

		// Create new user
		result, err := util.DB.Exec(
			"INSERT INTO tblUsers(username, email, auth_provider) VALUES(?, ?, ?)",
			user.Name, user.Email, "google",
		)
		if err != nil {
			log.Printf("User creation failed: %v", err)
			http.Redirect(w, r, "/sign-in?error=user_creation_failed", http.StatusTemporaryRedirect)
			return
		}

		id, _ := result.LastInsertId()
		userID = int(id)
		isNewUser = true
	} else if err != nil {
		log.Printf("Database error: %v", err)
		http.Redirect(w, r, "/sign-in?error=database_error", http.StatusTemporaryRedirect)
		return
	}

	// Create session token
	sessionToken := handler.CreateSession()

	// Delete any existing sessions for this user
	if userID != 0 {
		handler.DeleteSession(userID)
	}
	err = repositories.DeleteSessionByUser(userID)
	if err != nil {
		log.Printf("Failed to delete session token: %v", err)
		http.Redirect(w, r, "/sign-in?error=session_error", http.StatusTemporaryRedirect)
		return
	}

	// Enable CORS
	handler.EnableCors(w)

	// Set session cookie and data
	handler.SetSessionCookie(w, sessionToken)
	handler.SetSessionData(sessionToken, "userId", userID)
	handler.SetSessionData(sessionToken, "userEmail", user.Email)

	// Store session with expiry time
	expiryTime := time.Now().Add(24 * time.Hour)
	err = repositories.StoreSession(userID, sessionToken, expiryTime)
	if err != nil {
		log.Printf("Failed to store session token: %v", err)
		http.Redirect(w, r, "/sign-in?error=session_error", http.StatusTemporaryRedirect)
		return
	}

	// Redirect based on whether this is a new user or not
	if isNewUser {
		http.Redirect(w, r, "/home?status=new_user", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/home?status=returning_user", http.StatusSeeOther)
	}
}

// exchangeGoogleToken exchanges the authorization code for an access token
func exchangeGoogleToken(code string) (string, error) {
	data := url.Values{
		"code":          {code},
		"client_id":     {util.GoogleClientID},
		"client_secret": {util.GoogleClientSecret},
		"redirect_uri":  {RedirectBaseURL + "/auth/google/callback"},
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
