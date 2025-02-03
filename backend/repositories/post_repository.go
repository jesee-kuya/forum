package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jesee-kuya/forum/backend/models"
)

func GetPosts(db *sql.DB) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.post_title, p.body, p.created_on
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
		SELECT p.id, p.user_id, u.username, p.post_title, p.body, p.created_on
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

		err := rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.PostTitle, &post.Body, &post.CreatedOn)
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

// FilterPosts - Fetch posts based on category or user
func FilterPosts(db *sql.DB, filterType, filterValue string) ([]models.Post, error) {
	var query string
	var rows *sql.Rows
	var err error

	switch filterType {
	case "category":
		query = "SELECT * FROM tblPosts WHERE post_category = ?"
		rows, err = db.Query(query, filterValue)
	case "user":
		query = "SELECT * FROM tblPosts WHERE user_id = ?"
		rows, err = db.Query(query, filterValue)
	default:
		return nil, fmt.Errorf("invalid filter type")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Body, &post.ParentID, &post.CreatedOn); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}
