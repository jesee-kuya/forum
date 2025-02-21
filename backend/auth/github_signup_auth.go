package auth

const (
	GithubAuthURL   = "https://github.com/login/oauth/authorize"
	GithubTokenURL  = "https://github.com/login/oauth/access_token"
	GithubUserURL   = "https://api.github.com/user"
	GithubEmailsURL = "https://api.github.com/user/emails"
)

var (
	GithubClientID     = "GITHUB_CLIENT_ID"
	GithubClientSecret = "GITHUB_CLIENT_SECRET"
)

type GitHubUser struct {
	Login, Email string
}
