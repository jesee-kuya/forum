package route

import (
	"net/http"

	"github.com/jesee-kuya/forum/backend/handler"
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
	r.HandleFunc("/signin", handler.LoginHandler)
	r.HandleFunc("/signup", handler.SignupHandler)
	r.HandleFunc("/upload", handler.UploadMedia)

	// r.HandleFunc("/", handler.GetAllPosts(db, tmpl))

	return r
}
