package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
)

func GetAllPosts(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := repositories.GetPosts(db)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set content type to html and render the template
		w.Header().Set("Content-Type", "text/html")

		err = tmpl.ExecuteTemplate(w, "index.html", struct {
			Posts []models.Post
		}{Posts: posts})

		if err = json.NewEncoder(w).Encode(posts); err != nil {
			log.Printf("Failed to encode posts to JSON: %v", err)
			ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
