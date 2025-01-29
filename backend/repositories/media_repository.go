package repositories

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveMedia saves the file and stores metadata in the database.
func SaveMedia(db *sql.DB, file multipart.File, filename string, uploaderID int) (string, error) {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Println("Failed to create uploads directory:", err)
		return "", err
	}

	ext := filepath.Ext(filename)
	allowedExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".mp4": true, ".avi": true, ".mov": true,
		".wav": true,
	}
	if !allowedExtensions[ext] {
		log.Println("Unsupported file format:", ext)
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}

	newFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), filename)
	filePath := fmt.Sprintf("uploads/%s", newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Failed to create file:", err)
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println("Failed to save file:", err)
		return "", err
	}

	query := `INSERT INTO tblMedia (filename, file_path, file_type, uploader_id) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, newFilename, filePath, ext, uploaderID)
	if err != nil {
		log.Println("Failed to store media metadata:", err)
		return "", err
	}

	return filePath, nil
}
