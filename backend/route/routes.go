package route

import (
	"net/http"

	"github.com/jesee-kuya/forum/backend/handler"
	"github.com/jesee-kuya/forum/backend/middleware"
	"github.com/jesee-kuya/forum/backend/util"
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
	r.Handle("/upload", middleware.SessionMiddleware(http.HandlerFunc(handler.CreatePost)))
	r.Handle("/logout", middleware.SessionMiddleware(http.HandlerFunc(handler.LogoutHandler)))
	r.Handle("/comments", middleware.SessionMiddleware(http.HandlerFunc(handler.CommentHandler)))
	r.Handle("/reaction", middleware.SessionMiddleware(http.HandlerFunc(handler.ReactionHandler)))
	r.Handle("/likes", middleware.SessionMiddleware(http.HandlerFunc(handler.ReactionHandler)))
	r.Handle("/dilikes", middleware.SessionMiddleware(http.HandlerFunc(handler.ReactionHandler)))
	r.Handle("/filter", middleware.SessionMiddleware(http.HandlerFunc(handler.FilterPosts)))
	r.HandleFunc("/api/posts", handler.GetAllPostsAPI(util.DB))
	r.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleGetPosts(w, r, util.DB)
	})
	return r
}
