package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jesee-kuya/forum/backend/models"
)

var db *sql.DB

func InitRepo(database *sql.DB) {
	db = database
}

func GetPosts(postType string) ([]models.Post, error) {
	query := "SELECT * FROM tblposts WHERE post_type = ?"
	rows, err := db.Query(query, postType)
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

func CreatePost(post models.Post) (int64, error) {
	query := "INSERT INTO tblPosts (user_id, title, body, parent_id, post_type) VALUES (?, ?, ?, ?, ?)"
	result, err := db.Exec(query, post.UserID, post.TITLE, post.Body, post.ParentID, post.PostType)
	if err != nil {
		return 0, fmt.Errorf("failed to insert post: %v", err)
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve post ID: %v", err)
	}

	return postID, nil
}

func GetPostsByCategory(categoryID int) ([]models.Post, error) {
	query := `
        SELECT p.* 
        FROM tblposts p 
        JOIN tblpostcategories pc ON p.id = pc.post_id 
        WHERE pc.category_id = ?`
	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts by category: %v", err)
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

func AddReaction(reaction models.Reaction) error {
	// we'll first check if the user already reacted to the post
	checkQuery := "SELECT id FROM tblReactions WHERE user_id = ? AND post_id = ?"
	row := db.QueryRow(checkQuery, reaction.UserID, reaction.PostID)

	var existingID int
	err := row.Scan(&existingID)
	if err == nil {
		return fmt.Errorf("user has already reacted to this post")
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing reactions: %v", err)
	}

	// else we'll insert new reaction
	query := "INSERT INTO tblReactions (reaction, user_id, post_id) VALUES (?, ?, ?)"
	_, err = db.Exec(query, reaction.Reaction, reaction.UserID, reaction.PostID)
	if err != nil {
		return fmt.Errorf("failed to add reaction: %v", err)
	}

	return nil
}

func GetPostsByUser(userID int) ([]models.Post, error) {
	query := "SELECT * FROM tblposts WHERE user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts by user: %v", err)
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

func GetLikedPostsByUser(userID int) ([]models.Post, error) {
	query := `
        SELECT p.* 
        FROM tblposts p 
        JOIN tblreactions r ON p.id = r.post_id 
        WHERE r.user_id = ? AND r.reaction = 'like'`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch liked posts: %v", err)
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

func GetPostsWithPagination(postType string, limit, offset int) ([]models.Post, error) {
	query := "SELECT * FROM tblposts WHERE post_type = ? LIMIT ? OFFSET ?"
	rows, err := db.Query(query, postType, limit, offset)
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

func GetPostByID(postID string) (models.Post, error) {
	query := "SELECT * FROM tblPosts WHERE id = ?"
	row := db.QueryRow(query, postID)

	var post models.Post
	err := row.Scan(&post.ID, &post.UserID, &post.Body, &post.ParentID, &post.CreatedOn, &post.PostType)
	if err == sql.ErrNoRows {
		return post, fmt.Errorf("post not found")
	} else if err != nil {
		return post, fmt.Errorf("failed to fetch post: %v", err)
	}

	return post, nil
}

// LikePost - Increments the like count for a post
func LikePost(postID, userID int) error {
	query := "INSERT INTO tblReactions (reaction, user_id, post_id) VALUES (?, ?, ?)"
	_, err := db.Exec(query, "like", userID, postID)
	if err != nil {
		return fmt.Errorf("failed to like post: %v", err)
	}
	return nil
}

// CreateComment - Adds a comment to a post
func CreateComment(userID int, body string, parentID int) error {
	query := "INSERT INTO tblPosts (user_id, body, parent_id, post_type) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, userID, body, parentID, "comment")
	if err != nil {
		return fmt.Errorf("failed to create comment: %v", err)
	}
	return nil
}

// FilterPosts - Fetch posts based on category or user
func FilterPosts(filterType, filterValue string) ([]models.Post, error) {
	var query string
	var rows *sql.Rows
	var err error

	switch filterType {
	case "category":
		query = "SELECT * FROM tblPosts WHERE post_type = ?"
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
		if err := rows.Scan(&post.ID, &post.UserID, &post.Body, &post.ParentID, &post.CreatedOn, &post.PostType); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}
