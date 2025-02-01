package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

func GetAllPosts(db *sql.DB, tmpl *template.Template, posts []models.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fetch comments for each post
		for i, post := range posts {
			comments, err := repositories.GetComments(db, post.ID)
			if err != nil {
				log.Printf("Failed to get comments: %v", err)
				util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			posts[i].Comments = comments
		}

		// Set content type to text/html and serve the index page
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		err := tmpl.ExecuteTemplate(w, "index.html", struct {
			Posts []models.Post
		}{Posts: posts})
		if err != nil {
			log.Printf("Failed to render template: %v", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func GetAllPostsAPI(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := repositories.GetPosts(db)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// fetch comments for each post
		for i, post := range posts {
			comments, err := repositories.GetComments(db, post.ID)
			if err != nil {
				log.Printf("Failed to get posts: %v", err)
				util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			posts[i].Comments = comments
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(posts); err != nil {
			log.Printf("Failed to encode posts to JSON: %v", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// FilterPosts - Handles filtering posts by category or user
func FilterPosts(w http.ResponseWriter, r *http.Request) {
	filterType := r.URL.Query().Get("type")
	filterValue := r.URL.Query().Get("value")

	posts, err := repositories.FilterPosts(util.DB, filterType, filterValue)
	if err != nil {
		log.Printf("Failed to filter posts: %v", err)
		util.ErrorHandler(w, "Could not filter posts", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Printf("Failed to encode posts: %v", err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Println("Invalid request method")
		return
	}

	var req struct {
		PostID int `json:"post_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.ErrorHandler(w, "Bad request", http.StatusBadRequest)
		log.Printf("Failed to decode response body: %v\n", err)
		return
	}

	log.Printf("Received post_id: %d\n", req.PostID)

	newLikeCount, err := repositories.UpdateLikeCount(util.DB, req.PostID, 3)
	if err != nil {
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		log.Printf("Failed updating likes count: %v\n", err)
		return
	}

	resp := map[string]int{"newLikeCount": newLikeCount}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
