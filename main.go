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

	// addPost(db)

	posts, err := repositories.GetPosts(db)
	if err != nil {
		fmt.Println("Failed to fetch posts", err)
		return
	}

	fmt.Println(posts)
}

func addPost(db *sql.DB) {
	defer db.Close()

	post := models.Post{
		UserID:       1,
		PostTitle:    "Donald Trump",
		Body:         "You are eithe a male or female",
		PostCategory: "Politics",
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
