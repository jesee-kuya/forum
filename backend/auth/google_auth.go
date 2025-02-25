package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
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

// GoogleAuth handles both sign-in and sign-up flows using a flow parameter
func GoogleAuth(w http.ResponseWriter, r *http.Request) {
	// Get the flow type from query parameter (signup or signin)
	flow := r.URL.Query().Get("flow")
	if flow != "signup" && flow != "signin" {
		flow = "signup" // Default to signup if not specified
	}

	// Generate a random state and set it as a cookie to prevent CSRF attacks
	state := generateStateCookie(w, flow)

	// Construct the Google OAuth 2.0 authorization URL with necessary parameters
	redirectURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid email profile&state=%s&prompt=select_account&access_type=offline",
		GoogleAuthURL,
		util.GoogleClientID,
		url.QueryEscape("http://localhost:9000/auth/google/callback"),
		state,
	)

	// Set CORS headers for popup
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:9000")

	// Redirect the user to Google's OAuth 2.0 server
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the callback from the Google OAuth 2.0 server for both signup and signin
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Validate the state to prevent CSRF attacks
	flow, err := validateState(r)
	if err != nil {
		log.Printf("State validation failed: %v", err)

		return
	}

	// Get the authorization code from the query parameter
	code := r.URL.Query().Get("code")
	token, err := exchangeGoogleToken(code)
	if err != nil {
		log.Printf("Token exchange failed: %v\n", err)

		return
	}

	// Get the user information from the Google UserInfo endpoint
	user, err := getGoogleUser(token)
	if err != nil {
		log.Printf("Failed to get user info: %v\n", err)

		return
	}

	// Handle based on the flow type
	switch flow {
	case "signup":
		// Attempt to create or authenticate the user
		if handleUserAuth(w, user.Email, user.Name) {
			log.Printf("User authentication successful")
		} else {
			log.Printf("User authentication failed")
		}

	case "signin":
		// Check if the user exists in the database
		var userID int
		err = util.DB.QueryRow("SELECT id FROM tblUsers WHERE email = ?", user.Email).Scan(&userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return
			}
			log.Printf("Database error: %v", err)

			return
		}

		sessionToken := handler.CreateSession()

		// Delete any existing sessions for this user
		if userID != 0 {
			handler.DeleteSession(userID)
		}
		err = repositories.DeleteSessionByUser(userID)
		if err != nil {
			log.Printf("Failed to delete session token: %v", err)

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

			return
		}

	}
}

// exchangeGoogleToken exchanges the authorization code for an access token from Google.
func exchangeGoogleToken(code string) (string, error) {
	// Prepare the data for the token exchange request
	data := url.Values{
		"code":          {code},
		"client_id":     {util.GoogleClientID},
		"client_secret": {util.GoogleClientSecret},
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

// getGoogleUser retrieves the user information from the Google UserInfo endpoint.
func getGoogleUser(token string) (*GoogleUser, error) {
	// Make a GET request to Google's UserInfo endpoint
	req, err := http.NewRequest("GET", GoogleUserInfo, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response to extract the user information
	var user GoogleUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// generateStateCookie generates a random state and sets it as a cookie. It includes the flow type
// (signup or signin) in the state to differentiate the flows in the callback.
func generateStateCookie(w http.ResponseWriter, flow string) string {
	// Generate random state
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Failed to generate random state: %v", err)
		return ""
	}
	state := base64.URLEncoding.EncodeToString(b)

	// Store both state and flow type
	stateData := fmt.Sprintf("%s:%s", state, flow)

	// Set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    stateData,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return state
}

// validateState validates the state parameter in the callback to prevent CSRF attacks.
// It returns the flow type (signup or signin) if the state is valid.
func validateState(r *http.Request) (string, error) {
	// Get the cookie
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		return "", errors.New("oauth state cookie not found")
	}

	// Split the state data to get original state and flow type
	stateParts := strings.Split(cookie.Value, ":")
	if len(stateParts) != 2 {
		return "", errors.New("invalid oauth state format")
	}
	originalState, flow := stateParts[0], stateParts[1]

	// Get the state from the request
	receivedState := r.URL.Query().Get("state")
	if receivedState == "" {
		return "", errors.New("missing state parameter")
	}

	// Compare the states
	if receivedState != originalState {
		return "", errors.New("state mismatch")
	}

	return flow, nil
}
