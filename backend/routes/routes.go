package routes

import (
	"net/http"

	handler "github.com/jesee-kuya/forum/backend/handler"
	"github.com/jesee-kuya/forum/backend/middleware"
)

func RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	// Public routes
	r.HandleFunc("/api/register", handler.RegisterUser)
	r.HandleFunc("/api/login", handler.Login)
	r.HandleFunc("/api/posts", handler.GetAllPosts)

	// Protected routes (authentication required)
	protectedRoutes := map[string]http.HandlerFunc{
		"/api/posts/create":  handler.CreatePost,
		"/api/posts/like":    handler.LikePost,
		"/api/posts/comment": handler.CreateComment,
		"/api/posts/filter":  handler.FilterPosts,
		// "/api/upload":        handler.UploadMedia,
	}

	for route, handler := range protectedRoutes {
		r.HandleFunc(route, middleware.Authenticate(handler))
	}

	return r
}
