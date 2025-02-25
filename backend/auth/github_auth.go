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
	GithubAuthURL   = "https://github.com/login/oauth/authorize"
	GithubTokenURL  = "https://github.com/login/oauth/access_token"
	GithubUserURL   = "https://api.github.com/user"
	GithubEmailsURL = "https://api.github.com/user/emails"
	RedirectBaseURL = "http://localhost:9000"
)

type GitHubUser struct {
	Login, Email string
}

// GitHubAuth handles both sign-in and sign-up flows using a flow parameter
func GitHubAuth(w http.ResponseWriter, r *http.Request) {
	// Get the flow type from query parameter (signup or signin)
	flow := r.URL.Query().Get("flow")
	if flow != "signup" && flow != "signin" {
		flow = "signup" // Default to signup if not specified
	}

	state := generateGithubStateCookie(w, flow)

	params := url.Values{
		"client_id":    {util.GithubClientID},
		"redirect_uri": {RedirectBaseURL + "/auth/github/callback"},
		"scope":        {"user:email"},
		"state":        {state},
		"prompt":       {"consent"},
	}

	redirectURL := fmt.Sprintf("%s?%s", GithubAuthURL, params.Encode())
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// GitHubCallback handles the callback from GitHub's OAuth 2.0 server for both signup and signin
func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		log.Printf("Cookie error: %v", err)
		return
	}

	// Split the state data to get original state and flow type
	stateParts := strings.Split(cookie.Value, ":")
	if len(stateParts) != 2 {
		return
	}
	originalState, flow := stateParts[0], stateParts[1]

	// Validate state from URL against original state
	if r.URL.Query().Get("state") != originalState {
		log.Printf("State mismatch")
		return
	}

	code := r.URL.Query().Get("code")
	token, err := exchangeGitHubToken(code, "/auth/github/callback")
	if err != nil {
		log.Printf("Token exchange failed: %v\n", err)
		return
	}

	user, err := getGitHubUser(token)
	if err != nil {
		log.Printf("Failed to get user info: %v\n", err)
		return
	}

	// Handle based on flow type
	switch flow {
	case "signup":
		if handleUserAuth(w, user.Email, user.Login) {
			log.Printf("User signup successful")
		} else {
			log.Printf("User signup failed")
		}

	case "signin":
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
		if userID != 0 {
			handler.DeleteSession(userID)
		}
		err = repositories.DeleteSessionByUser(userID)
		if err != nil {
			log.Printf("Failed to delete session token: %v", err)
			return
		}

		handler.EnableCors(w)
		handler.SetSessionCookie(w, sessionToken)
		handler.SetSessionData(sessionToken, "userId", userID)
		handler.SetSessionData(sessionToken, "userEmail", user.Email)

		expiryTime := time.Now().Add(24 * time.Hour)
		err = repositories.StoreSession(userID, sessionToken, expiryTime)
		if err != nil {
			log.Printf("Failed to store session token: %v", err)
			return
		}
	}
}

// exchangeGitHubToken exchanges the authorization code for an access token from GitHub.
func exchangeGitHubToken(code, redirectPath string) (string, error) {
	data := url.Values{
		"code":          {code},
		"client_id":     {util.GithubClientID},
		"client_secret": {util.GithubClientSecret},
		"redirect_uri":  {RedirectBaseURL + redirectPath},
	}

	req, _ := http.NewRequest("POST", GithubTokenURL, nil)
	req.URL.RawQuery = data.Encode()
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("%s: %s", result.Error, result.ErrorDesc)
	}

	return result.AccessToken, nil
}

// getGitHubUser retrieves the GitHub user profile and email (if not public) from the GitHub API given an access token.
func getGitHubUser(token string) (*GitHubUser, error) {
	req, _ := http.NewRequest("GET", GithubUserURL, nil)
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	// Fetch email from emails endpoint if not provided/not public
	if user.Email == "" {
		req, _ = http.NewRequest("GET", GithubEmailsURL, nil)
		req.Header.Set("Authorization", "token "+token)
		req.Header.Set("Accept", "application/json")

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
			return nil, err
		}

		// Find primary email
		for _, email := range emails {
			if email.Primary && email.Verified {
				user.Email = email.Email
				break
			}
		}
	}

	return &user, nil
}

/*
generateGithubStateCookie generates a random state and sets it as a cookie. It returns the generated state string. The state string will be in the format "<state>:<flow_type>" where flow_type is either "signup" or "signin".
*/
func generateGithubStateCookie(w http.ResponseWriter, flowType string) string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Printf("error generating random state: %v", err)
		return ""
	}

	state := base64.URLEncoding.EncodeToString(b)

	// Store both state and flow type
	stateData := fmt.Sprintf("%s:%s", state, flowType)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    stateData,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode,
	})
	return state
}
