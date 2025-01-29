package models

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// User model
type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	JoinedOn time.Time `json:"joined_on"`
}

// Post model
type Post struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	PostTitle    string    `json:"post_title"`
	Body         string    `json:"body"`
	ParentID     int       `json:"parent_id"`
	CreatedOn    time.Time `json:"created_on"`
	PostCategory string    `json:"post_category"`
	PostStatus   string    `json:"post_status"`
}

// File model
type File struct {
	ID       int    `json:"id"`
	PostID   int    `json:"post_id"`
	FileName string `json:"file_path"`
	FileType string `json:"file_type"`
}

// Reaction model
type Reaction struct {
	ID       int    `json:"id"`
	Reaction string `json:"reaction"`
	UserID   int    `json:"user_id"`
	PostID   int    `json:"post_id"`
}
