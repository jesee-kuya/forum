package models

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// User model
type User struct {
	ID                int       `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          string    `json:"password,omitempty"`
	ConfirmedPassword string    `json:"confirmed-password"`
	JoinedOn          time.Time `json:"joined_on"`
}

// Post model
type Post struct {
	ID           int        `json:"id"`
	UserID       int        `json:"user_id"`
	UserName     string     `json:"username"`
	PostTitle    string     `json:"post_title"`
	Body         string     `json:"body"`
	ParentID     *int       `json:"parent_id"`
	CreatedOn    time.Time  `json:"created_on"`
	PostStatus   string     `json:"post_status"`
	Likes        int        `json:"likes"`
	Dislikes     int        `json:"dislikes"`
	CommentCount int        `json:"comment_count"`
	Categories   []Category `json:"categorie"`
	MediaURL     string     `json:"imageurl"`
	Comments     []Post     `json:"comments"`
}

// Category model
type Category struct {
	ID           int    `json:"id"`
	PostID       int    `json:"post_id"`
	CategoryName string `json:"category"`
}

// Reaction model
type Reaction struct {
	ID             int    `json:"id"`
	Reaction       string `json:"reaction"`
	ReactionStatus string `json:"reaction_status"`
	UserID         int    `json:"user_id"`
	PostID         int    `json:"post_id"`
}
