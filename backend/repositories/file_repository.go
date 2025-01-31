package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jesee-kuya/forum/backend/models"
)

func GetMediaFiles(db *sql.DB, id int) ([]models.MediaFile, error) {
	query := `
		SELECT * FROM tblMediaFiles 
		WHERE post_id = ? AND file_status = 'visible'
	`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var files []models.MediaFile

	for rows.Next() {
		file := models.MediaFile{}

		err := rows.Scan(&file.ID, &file.PostID, &file.FileName, &file.FileType, &file.FileStatus)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		files = append(files, file)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return files, nil
}
