package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jesee-kuya/forum/backend/models"
)

func GetPosts(db *sql.DB) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.post_title, p.body, p.created_on, p.post_category
		FROM tblPosts p
		JOIN tblUsers u ON p.user_id = u.id
		WHERE p.parent_id IS NULL AND p.post_status = 'visible'`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	posts, err := processSQLData(rows)

	return posts, err
}

func GetComments(db *sql.DB, id int) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.post_title, p.body, p.created_on, p.post_category
		FROM tblPosts p
		JOIN tblUsers u ON p.user_id = u.id
		WHERE p.parent_id = ? AND p.post_status = 'visible'
	`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	posts, err := processSQLData(rows)

	return posts, err
}

func processSQLData(rows *sql.Rows) ([]models.Post, error) {
	var posts []models.Post

	for rows.Next() {
		post := models.Post{}

		err := rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.PostTitle, &post.Body, &post.CreatedOn, &post.PostCategory)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		posts = append(posts, post)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return posts, nil
}
