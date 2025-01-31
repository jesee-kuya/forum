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
