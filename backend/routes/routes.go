package routes

import (
	"net/http"

	"github.com/jesee-kuya/forum/backend/controllers"
	"github.com/jesee-kuya/forum/backend/middleware"
)

func RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	// Public routes
	r.HandleFunc("/api/register", controllers.RegisterUser)
	r.HandleFunc("/api/login", controllers.Login)
	r.HandleFunc("/api/posts", controllers.GetAllPosts)

	// Protected routes (authentication required)
	protectedRoutes := map[string]http.HandlerFunc{
		"/api/posts/create":  controllers.CreatePost,
		"/api/posts/like":    controllers.LikePost,
		"/api/posts/comment": controllers.CreateComment,
		"/api/posts/filter":  controllers.FilterPosts,
		// "/api/upload":        controllers.UploadMedia,
	}

	for route, handler := range protectedRoutes {
		r.HandleFunc(route, middleware.Authenticate(handler))
	}

	return r
}
