package route

import (
	"net/http"
	"time"

	"github.com/jesee-kuya/forum/backend/handler"
	"github.com/jesee-kuya/forum/backend/middleware"
	openauth "github.com/jesee-kuya/forum/backend/open_auth"
)

func InitRoutes() *http.ServeMux {
	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("./frontend"))
	r.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	uploadFs := http.FileServer(http.Dir("./uploads"))
	r.Handle("/uploads/", http.StripPrefix("/uploads/", uploadFs))

	// App routes
	r.HandleFunc("/home", middleware.Authenticate(middleware.RateLimiter(handler.IndexHandler, 50, time.Second)))
	r.HandleFunc("/", middleware.RateLimiter(handler.HomeHandler, 50, time.Second))
	r.HandleFunc("/sign-in", middleware.RateLimiter(handler.LoginHandler, 50, time.Second))
	r.HandleFunc("/sign-up", middleware.RateLimiter(handler.SignupHandler, 50, time.Second))
	r.HandleFunc("/upload", middleware.Authenticate(middleware.RateLimiter(handler.CreatePost, 50, time.Second)))
	r.HandleFunc("/logout", middleware.Authenticate(middleware.RateLimiter(handler.LogoutHandler, 50, time.Second)))
	r.HandleFunc("/comments", middleware.Authenticate(middleware.RateLimiter(handler.CommentHandler, 50, time.Second)))
	r.HandleFunc("/reaction", middleware.Authenticate(middleware.RateLimiter(handler.ReactionHandler, 50, time.Second)))
	r.HandleFunc("/likes", middleware.Authenticate(middleware.RateLimiter(handler.ReactionHandler, 50, time.Second)))
	r.HandleFunc("/dilikes", middleware.Authenticate(middleware.RateLimiter(handler.ReactionHandler, 50, time.Second)))
	r.HandleFunc("/filter", middleware.RateLimiter(handler.FilterPosts, 50, time.Second))

	r.HandleFunc("/validate", handler.ValidateInputHandler)


	r.HandleFunc("/auth/google", openauth.GoogleAuth)
	r.HandleFunc("/auth/google/callback", openauth.GoogleCallback)

	r.HandleFunc("/auth/github", openauth.GitHubAuth)
	r.HandleFunc("/auth/github/callback", openauth.GitHubCallback)
	return r
}
