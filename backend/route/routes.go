package route

import (
	"net/http"

	"github.com/jesee-kuya/forum/backend/handler"
	"github.com/jesee-kuya/forum/backend/util"
)

func InitRoutes() *http.ServeMux {
	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("./frontend"))
	r.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	uploadFs := http.FileServer(http.Dir("./uploads"))
	r.Handle("/uploads/", http.StripPrefix("/uploads/", uploadFs))

	// App routes
	r.HandleFunc("/home", handler.IndexHandler)
	r.HandleFunc("/", handler.HomeHandler)
	r.HandleFunc("/sign-in", handler.LoginHandler)
	r.HandleFunc("/sign-up", handler.SignupHandler)
	r.HandleFunc("/upload", handler.CreatePost)
	r.HandleFunc("/logout", handler.LogoutHandler)
	r.HandleFunc("/comments", handler.CommentHandler)
	r.HandleFunc("/likes", handler.LikeHandler)
	r.HandleFunc("/dilikes", handler.DislikeHandler)

	r.HandleFunc("/api/posts", handler.GetAllPostsAPI(util.DB))
	return r
}
