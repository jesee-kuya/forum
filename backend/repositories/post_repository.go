package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jesee-kuya/forum/backend/models"
)

func GetPosts(db *sql.DB) ([]models.Post, error) {
	query := `
		SELECT id, user_id, post_title, body, created_on, post_category
		FROM tblPosts
		WHERE parent_id IS NULL AND post_status = 'visible';
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		post := models.Post{}

		err := rows.Scan(&post.ID, &post.UserID, &post.PostTitle, &post.Body, &post.CreatedOn, &post.PostCategory)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		posts = append(posts, post)
	}

	// Check for errors after iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return posts, nil
}
