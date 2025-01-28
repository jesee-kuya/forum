package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jesee-kuya/forum/backend/models"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

// insertPost inserts a Post into the tblPosts table
func insertFile(db *sql.DB, file models.File) (int64, error) {
	query := `
		INSERT INTO tblFiles (post_id, file_name, file_type)
		VALUES (?, ?, ?)
	`

	result, err := db.Exec(query, file.PostID, file.FileName, file.FileType)
	if err != nil {
		return 0, fmt.Errorf("failed to insert post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return id, nil
}

// deletePost deletes a record from tblPosts based on its ID
func deleteFile(db *sql.DB, fileID int) error {
	query := `
		DELETE FROM tblFiles
		WHERE ID = ?;
	`

	// Execute the DELETE statement
	result, err := db.Exec(query, fileID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	// Check if a row was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no post found with ID %d", fileID)
	}

	log.Printf("Successfully deleted post with ID: %d", fileID)
	return nil
}

func RemoveFile(db *sql.DB) {

	defer db.Close()

	id := 1

	// Delete the post
	err := deleteFile(db, id)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	log.Println("Post deletion complete.")
}

// func AddFile(db *sql.DB) {

// 	defer db.Close()

// 	// Example File data
// 	file := models.File{
// 		PostID:   1,
// 		FileName: "example.jpg",
// 		FileType: "Profile Image",
// 	}

// 	// Insert the post into the database
// 	id, err := insertPost(db, file)
// 	if err != nil {
// 		log.Fatalf("failed to insert post: %v", err)
// 	}

// 	log.Printf("Post inserted successfully with ID: %d", id)
// }
