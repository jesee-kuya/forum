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
	//addFile(db)
	deleteFile(db)
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
