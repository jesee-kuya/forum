package auth

const (
	GithubAuthURL   = "https://github.com/login/oauth/authorize"
	GithubTokenURL  = "https://github.com/login/oauth/access_token"
	GithubUserURL   = "https://api.github.com/user"
	GithubEmailsURL = "https://api.github.com/user/emails"
)

type GitHubUser struct {
	Login, Email string
}
