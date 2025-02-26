package openauth

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
)

type GitHubUser struct {
	Login, Email string
}

// GitHubAuth initiates the GitHub authentication process (signup or signin)
func GitHubAuth(w http.ResponseWriter, r *http.Request) {
	state := generateStateCookie(w)

	params := url.Values{
		"client_id":    {util.GithubClientID},
		"redirect_uri": {RedirectBaseURL + "/auth/github/callback"},
		"scope":        {"user:email"},
		"state":        {state},
		"prompt":       {"consent"},
	}

	redirectURL := fmt.Sprintf("%s?%s", GithubAuthURL, params.Encode())
	w.Header().Set("Access-Control-Allow-Origin", RedirectBaseURL)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// GitHubCallback handles the callback from GitHub's OAuth server
func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	if err := validateState(r); err != nil {
		log.Printf("State validation failed: %v", err)
		http.Redirect(w, r, "/auth-error?type=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := exchangeGitHubToken(code)
	if err != nil {
		log.Printf("Token exchange failed: %v\n", err)
		http.Redirect(w, r, "/auth-error?type=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	user, err := getGitHubUser(token)
	if err != nil {
		log.Printf("Failed to get user info: %v\n", err)
		http.Redirect(w, r, "/auth-error?type=user_info_failed", http.StatusTemporaryRedirect)
		return
	}

	var userID int
	err = util.DB.QueryRow("SELECT id FROM tblUsers WHERE email = ?", user.Email).Scan(&userID)

	isNewUser := false

	// If user doesn't exist, create a new one
	if errors.Is(err, sql.ErrNoRows) {
		var count int
		err = util.DB.QueryRow("SELECT COUNT(*) FROM tblUsers WHERE username = ?", user.Login).Scan(&count)
		if err != nil {
			log.Printf("Database error checking username: %v", err)
			http.Redirect(w, r, "/auth-error?type=database_error", http.StatusTemporaryRedirect)
			return
		}

		if count > 0 {
			// Username is taken, generate a unique one
			b := make([]byte, 3)
			rand.Read(b)
			suffix := base64.URLEncoding.EncodeToString(b)
			user.Login = fmt.Sprintf("%s_%s", user.Login, suffix)
		}

		// Create new user
		result, err := util.DB.Exec(
			"INSERT INTO tblUsers(username, email) VALUES(?, ?)",
			user.Login, user.Email,
		)
		if err != nil {
			log.Printf("User creation failed: %v", err)
			http.Redirect(w, r, "/auth-error?type=user_creation_failed", http.StatusTemporaryRedirect)
			return
		}

		id, _ := result.LastInsertId()
		userID = int(id)
		isNewUser = true
	} else if err != nil {
		log.Printf("Database error: %v", err)
		http.Redirect(w, r, "/auth-error?type=database_error", http.StatusTemporaryRedirect)
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
		http.Redirect(w, r, "/auth-error?type=session_error", http.StatusTemporaryRedirect)
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
		http.Redirect(w, r, "/auth-error?type=session_error", http.StatusTemporaryRedirect)
		return
	}

	// Redirect based on whether this is a new user or not
	if isNewUser {
		http.Redirect(w, r, "/home?status=new_user", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/home?status=returning_user", http.StatusSeeOther)
	}
}

// exchangeGitHubToken exchanges the authorization code for an access token
func exchangeGitHubToken(code string) (string, error) {
	data := url.Values{
		"code":          {code},
		"client_id":     {util.GithubClientID},
		"client_secret": {util.GithubClientSecret},
		"redirect_uri":  {RedirectBaseURL + "/auth/github/callback"},
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

// getGitHubUser retrieves the GitHub user profile and email
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
