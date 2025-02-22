package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jesee-kuya/forum/backend/util"
)

// GitHubSignIn initiates the GitHub sign-in process by redirecting the user to GitHub's OAuth 2.0 server for authentication.
func GitHubSignIn(w http.ResponseWriter, r *http.Request) {
	state := generateGithubStateCookie(w, "signin")

	params := url.Values{
		"client_id":    {util.GithubClientID},
		"redirect_uri": {RedirectBaseURL + "/auth/github/callback"},
		"scope":        {"user:email"},
		"state":        {state},
		"prompt":       {"consent select_account"},
	}

	redirectURL := fmt.Sprintf("%s?%s", GithubAuthURL, params.Encode())
	w.Header().Set("Access-Control-Allow-Origin", RedirectBaseURL)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
