package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/jesee-kuya/forum/backend/database"
	"github.com/jesee-kuya/forum/backend/repositories"
)

func GetAllPosts(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := database.CreateConnection()
		posts, err := repositories.GetPosts(db)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			ErrorHandler(w, "Internal Server Error", http.StatusMethodNotAllowed)
			return
		}

		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Printf("Failed to encode posts to JSON: %v", err)
			ErrorHandler(w, "Internal Server Error", http.StatusMethodNotAllowed)
			return
		}
	}
}
