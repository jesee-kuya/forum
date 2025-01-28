package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jesee-kuya/forum/backend/models"
)

var db *sql.DB

func GetPosts() ([]models.Post, error) {
	query := "SELECT * FROM tblposts"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	posts := make([]models.Post, 0)

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Body, &post.ParentID, &post.CreatedOn, &post.PostType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}
