package main

import (
	"database/sql"
	"fmt"

	"github.com/jesee-kuya/forum/backend/database"
	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
)

func main() {
	db := database.CreateConnection()
	// addFile(db)
	// deleteFile(db)
	// deletePost(db)

	//addPost(db)
	//addComment(db)

	posts, err := repositories.GetPosts(db)
	if err != nil {
		fmt.Println("Failed to fetch posts", err)
		return
	}

	fmt.Println(posts)
}

func deletePost(db *sql.DB) {
	defer db.Close()

	repositories.DeleteRecord(db, "tblPosts", "post_status", 1)
}

func addComment(db *sql.DB) {
	defer db.Close()

	post := models.Post{
		UserID:       6,
		PostTitle:    "Deepseek AI",
		Body:         "Make China Great Again",
		ParentID:     4,
		PostCategory: "Technology",
	}

	id, err := repositories.InsertRecord(db, "tblPosts", []string{"user_id", "post_title", "body", "parent_id", "post_category"}, post.UserID, post.PostTitle, post.Body, post.ParentID, post.PostCategory)
	if err != nil {
		fmt.Println("Failed to add the post to the database", err)
		return
	}

	fmt.Println(id)
}

func addPost(db *sql.DB) {
	defer db.Close()

	post := models.Post{
		UserID:       3,
		PostTitle:    "Deepseek AI",
		Body:         `It’s a side project. I call it DeepSeek.`,
		PostCategory: "Technology",
	}

	id, err := repositories.InsertRecord(db, "tblPosts", []string{"user_id", "post_title", "body", "post_category"}, post.UserID, post.PostTitle, post.Body, post.PostCategory)
	if err != nil {
		fmt.Println("Failed to add the post to the database", err)
		return
	}

	fmt.Println(id)
}

func addUser(db sql.DB) {
}

func addFile(db *sql.DB) {
	defer db.Close()
	file := models.File{
		PostID:   1,
		FileName: "example1.jpg",
		FileType: "upload",
	}

	id, err := repositories.InsertRecord(db, "tblFiles", []string{"post_id", "file_name", "file_type"}, file.PostID, file.FileName, file.FileType)
	if err != nil {
		fmt.Println("failed to add Record")
	}

	fmt.Println(id)
}

func deleteFile(db *sql.DB) {
	defer db.Close()

	repositories.DeleteRecord(db, "tblFiles", "file_status", 1)
}
