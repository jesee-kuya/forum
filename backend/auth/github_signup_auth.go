package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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

// GitHubSignUp initiates the GitHub sign-up process by redirecting the user to GitHub's OAuth 2.0 server for authentication.
func GitHubSignUp(w http.ResponseWriter, r *http.Request) {
	state := generateGithubStateCookie(w, "signup")

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

// GitHubCallback handles the callback from GitHub's OAuth 2.0 server after the user has granted the necessary permissions.
func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		log.Printf("Cookie error: %v", err)
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// Split the state data to get original state and flow type
	stateParts := strings.Split(cookie.Value, ":")
	if len(stateParts) != 2 {
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}
	originalState, flowType := stateParts[0], stateParts[1]

	// Validate state from URL against original state
	if r.URL.Query().Get("state") != originalState {
		log.Printf("State mismatch")
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := exchangeGitHubToken(code, "/auth/github/callback")
	if err != nil {
		log.Printf("Token exchange failed: %v\n", err)
		http.Redirect(w, r, "/sign-in?error=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	user, err := getGitHubUser(token)
	if err != nil {
		log.Printf("Failed to get user info: %v\n", err)
		http.Redirect(w, r, "/sign-in?error=user_info_failed", http.StatusTemporaryRedirect)
		return
	}

	// Handle based on flow type
	switch flowType {
	case "signup":
		if handleUserAuth(w, user.Email, user.Login) {
			log.Printf("User signup successful")
			// http.Redirect(w, r, "/sign-in", http.StatusTemporaryRedirect)
			http.Redirect(w, r, "/sign-in?status=success", http.StatusTemporaryRedirect)
		} else {
			log.Printf("User signup failed")
			http.Redirect(w, r, "/sign-in?error=auth_failed", http.StatusTemporaryRedirect)
		}

	case "signin":
		var userID int
		err = util.DB.QueryRow("SELECT id FROM tblUsers WHERE email = ?", user.Email).Scan(&userID)
		if err != nil {
			http.Redirect(w, r, "/sign-up?error=no_account", http.StatusTemporaryRedirect)
			return
		}

		sessionToken := handler.CreateSession()
		if userID != 0 {
			handler.DeleteSession(userID)
		}

		handler.EnableCors(w)
		handler.SetSessionCookie(w, sessionToken)
		handler.SetSessionData(sessionToken, "userId", userID)
		handler.SetSessionData(sessionToken, "userEmail", user.Email)

		expiryTime := time.Now().Add(24 * time.Hour)
		err = repositories.StoreSession(userID, sessionToken, expiryTime)
		if err != nil {
			log.Printf("Failed to store session token: %v", err)
			http.Redirect(w, r, "/sign-in?error=session_error", http.StatusTemporaryRedirect)
			return
		}

		// http.Redirect(w, r, "/home", http.StatusSeeOther)
		http.Redirect(w, r, "/home?status=success", http.StatusSeeOther)

	default:
		http.Redirect(w, r, "/sign-in?error=invalid_flow", http.StatusTemporaryRedirect)
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
