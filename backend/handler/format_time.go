package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

/*
FormatTimestamp converts the timestamp in the database to UTC format.
*/
func FormatTimestamp(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, user_id, username, post_title, body, created_on, media_url FROM posts")
	if err != nil {
		log.Printf("Failed fetching from database: %v\n", err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts, err := repositories.ProcessSQLData(rows)
	if err != nil {
		log.Printf("Failed processing database rows: %v\n", err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}

	// Convert timestamps to UTC format
	for i := range posts {
		posts[i].CreatedOn = posts[i].CreatedOn.UTC()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func HandleGetPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	posts, err := repositories.GetPosts(db)
	if err != nil {
		log.Println("error getting posts:", err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}

	for i := range posts {
		posts[i].CreatedOn = posts[i].CreatedOn.UTC()
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}
}
