package route

import (
	"net/http"

	"github.com/jesee-kuya/forum/backend/handler"
	"github.com/jesee-kuya/forum/backend/middleware"
	"github.com/jesee-kuya/forum/backend/util"
)

func InitRoutes() *http.ServeMux {
	r := http.NewServeMux()

	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("./frontend"))
	r.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	// Serve uploaded media files
	uploadFs := http.FileServer(http.Dir("./uploads"))
	r.Handle("/uploads/", http.StripPrefix("/uploads/", uploadFs))

	// App routes
	r.HandleFunc("/", handler.IndexHandler)
	r.HandleFunc("/sign-in", handler.LoginHandler)
	r.HandleFunc("/sign-up", handler.SignupHandler)
	r.HandleFunc("/upload", middleware.Authenticate(handler.CreatePost))

	// r.HandleFunc("/posts", handler.GetAllPosts(db, tmpl))
	r.HandleFunc("/api/posts", handler.GetAllPostsAPI(util.DB))

	return r
}
