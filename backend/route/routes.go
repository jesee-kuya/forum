package route

import (
	"net/http"

	"github.com/jesee-kuya/forum/backend/auth"
	"github.com/jesee-kuya/forum/backend/handler"
	"github.com/jesee-kuya/forum/backend/middleware"
)

func InitRoutes() *http.ServeMux {
	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("./frontend"))
	r.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	uploadFs := http.FileServer(http.Dir("./uploads"))
	r.Handle("/uploads/", http.StripPrefix("/uploads/", uploadFs))

	// App routes
	r.HandleFunc("/home", middleware.Authenticate(handler.IndexHandler))
	r.HandleFunc("/", handler.HomeHandler)
	r.HandleFunc("/sign-in", handler.LoginHandler)
	r.HandleFunc("/sign-up", handler.SignupHandler)
	r.HandleFunc("/upload", middleware.Authenticate(handler.CreatePost))
	r.HandleFunc("/logout", middleware.Authenticate(handler.LogoutHandler))
	r.HandleFunc("/comments", middleware.Authenticate(handler.CommentHandler))
	r.HandleFunc("/reaction", middleware.Authenticate(handler.ReactionHandler))
	r.HandleFunc("/likes", middleware.Authenticate(handler.ReactionHandler))
	r.HandleFunc("/dilikes", middleware.Authenticate(handler.ReactionHandler))
	r.HandleFunc("/filter", handler.FilterPosts)

	r.HandleFunc("/validate", handler.ValidateInputHandler)

	r.HandleFunc("/auth/google", auth.GoogleSignUp)
	r.HandleFunc("/auth/google/callback", auth.GoogleCallback)
	r.HandleFunc("/auth/google/signin", auth.GoogleSignIn)
	r.HandleFunc("/auth/google/signin/callback", auth.GoogleSignInCallback)

	// r.HandleFunc("/auth/github", auth.GitHubSignUp)
	// r.HandleFunc("/auth/github/callback", auth.GitHubCallback)
	// r.HandleFunc("/auth/google/signin", GitHubSignIn)
	// r.HandleFunc("/auth/google/signin/callback", GitHubSignInCallback)
	return r
}
