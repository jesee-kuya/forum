package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/database"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	db := database.CreateConnection()
	posts, err := repositories.GetPosts(db)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusMethodNotAllowed)
		return
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Printf("Failed to encode posts to JSON: %v", err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusMethodNotAllowed)
		return
	}
}

// FilterPosts - Handles filtering posts by category or user
func FilterPosts(w http.ResponseWriter, r *http.Request) {
	db := database.CreateConnection()
	filterType := r.URL.Query().Get("type")
	filterValue := r.URL.Query().Get("value")

	posts, err := repositories.FilterPosts(db, filterType, filterValue)
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
