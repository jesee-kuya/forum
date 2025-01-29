package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/handler"
	"github.com/jesee-kuya/forum/backend/repositories"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := repositories.GetPosts()
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		handler.ErrorHandler(w, "Internal Server Error", http.StatusMethodNotAllowed)
		return
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Printf("Failed to encode posts to JSON: %v", err)
		handler.ErrorHandler(w, "Internal Server Error", http.StatusMethodNotAllowed)
		return
	}
}
