package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jesee-kuya/forum/backend/models"
)

func GetCategories(db *sql.DB, id int) ([]models.Category, error) {
	query := `
		SELECT * FROM tblPostCategories
		WHERE post_id = ? 
	`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		category := models.Category{}

		err := rows.Scan(&category.ID, &category.PostID, &category.CategoryName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		categories = append(categories, category)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return categories, nil
}
