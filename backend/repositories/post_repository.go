package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jesee-kuya/forum/backend/models"
)

var PostQuery string

func GetPosts(db *sql.DB) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.post_title, p.body, p.created_on, p.media_url
		FROM tblPosts p
		JOIN tblUsers u ON p.user_id = u.id
		WHERE p.parent_id IS NULL AND p.post_status = 'visible'
		ORDER BY p.created_on DESC
		`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	posts, err := processSQLData(rows)
	if err != nil {
		return nil, fmt.Errorf("failed process posts: %v", err)
	}

	return posts, err
}

func GetComments(db *sql.DB, id int) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.post_title, p.body, p.created_on, p.media_url
		FROM tblPosts p
		JOIN tblUsers u ON p.user_id = u.id
		WHERE p.parent_id = ? AND p.post_status = 'visible'
		ORDER BY p.created_on DESC
	`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	posts, err := processSQLData(rows)
	if err != nil {
		return nil, fmt.Errorf("failed process comments: %v", err)
	}

	return posts, err
}

func FilterPostsByCategories(db *sql.DB, categories []string) ([]models.Post, error) {
	placeholders := strings.Repeat("?,", len(categories)-1) + "?"

	query := fmt.Sprintf(`
		SELECT DISTINCT p.id, p.user_id, u.username, p.post_title, p.body, p.created_on, p.media_url
		FROM tblPosts p
		JOIN tblUsers u ON p.user_id = u.id
		LEFT JOIN tblPostCategories c ON p.id = c.post_id 
		WHERE p.parent_id IS NULL 
		AND p.post_status = 'visible' 
		AND c.category IN (%s) ORDER BY p.created_on DESC`, placeholders)

	args := make([]interface{}, len(categories))
	for i, v := range categories {
		args[i] = v
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	posts, err := processSQLData(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to process posts: %v", err)
	}

	return posts, nil
}

func FilterPostsByUser(db *sql.DB, id int) ([]models.Post, error) {
	query := `
		SELECT DISTINCT p.id, p.user_id, u.username, p.post_title, p.body, p.created_on, p.media_url
		FROM tblPosts p
		JOIN tblUsers u ON p.user_id = u.id
		WHERE p.parent_id IS NULL AND p.post_status = 'visible' AND u.id = ?
		ORDER BY p.created_on DESC
		`

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	posts, err := processSQLData(rows)
	if err != nil {
		return nil, fmt.Errorf("failed process posts: %v", err)
	}

	return posts, err
}

// FilterPosts - Fetch posts based on category or user
func FilterPostsByLikes(db *sql.DB, id int) ([]models.Post, error) {
	query := `
		SELECT DISTINCT p.id, p.user_id, u.username, p.post_title, p.body, p.created_on, p.media_url
		FROM tblPosts p
		JOIN tblUsers u ON p.user_id = u.id
		LEFT JOIN tblReactions r ON p.id = r.post_id 
		WHERE p.parent_id IS NULL 
		AND p.post_status = 'visible' 
		AND r.reaction_status = 'clicked' 
		AND r.reaction = 'Like' 
		AND r.user_id = ?
		ORDER BY p.created_on DESC
		`

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	posts, err := processSQLData(rows)
	if err != nil {
		return nil, fmt.Errorf("failed process posts: %v", err)
	}

	return posts, err
}

func processSQLData(rows *sql.Rows) ([]models.Post, error) {
	var posts []models.Post

	for rows.Next() {
		post := models.Post{}

		err := rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.PostTitle, &post.Body, &post.CreatedOn, &post.MediaURL)
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
