package route

import (
	"net/http"

	"github.com/jesee-kuya/forum/backend/handler"
)

func InitRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/", handler.IndexHandler)
	r.HandleFunc("/signin", handler.LoginHandler)
	r.HandleFunc("/signup", handler.SignupHandler)
	r.HandleFunc("/upload", handler.UploadMedia)
	r.HandleFunc("/posts", handler.GetAllPosts)
	return r
}
