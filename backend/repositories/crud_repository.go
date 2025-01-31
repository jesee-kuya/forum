package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

// insertPost inserts a Post into the tblPosts table
func InsertRecord(db *sql.DB, table string, columns []string, values ...interface{}) (int64, error) {
	// Constructing column names and placeholders
	columnsStr := strings.Join(columns, ", ")
	placeholders := strings.Repeat("?, ", len(columns))
	placeholders = strings.TrimSuffix(placeholders, ", ")

	// Constructing the SQL query dynamically
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, columnsStr, placeholders)

	// Executing the query
	result, err := db.Exec(query, values...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert into %s: %w", table, err)
	}

	// Retrieving the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID for %s: %w", table, err)
	}

	return id, nil
}

// deletePost deletes a record from tblPosts based on its ID
func DeleteRecord(db *sql.DB, table, column string, id int) error {
	// Use a parameterized query for safety
	query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", table, column)

	// Execute the query safely with parameters
	result, err := db.Exec(query, "Deleted", id)
	if err != nil {
		return fmt.Errorf("failed to delete record from %s: %w", table, err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no record found with ID %d in %s", id, table)
	}

	log.Printf("Successfully marked record with ID %d as deleted in table %s", id, table)
	return nil
}
