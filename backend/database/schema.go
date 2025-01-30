package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

// executeSQLFile reads and executes SQL statements from a file on the given database connection
func executeSQLFile(db *sql.DB) error {
	content, err := os.ReadFile("./backend/database/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	// Split SQL statements and execute each one
	statements := string(content)
	_, err = db.Exec(statements)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statements: %w", err)
	}

	return nil
}

func CreateConnection() {
	// Open SQLite database connection
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Execute the SQL file
	err = executeSQLFile(db)
	if err != nil {
		log.Fatalf("failed to execute SQL file: %v", err)
	}

	log.Println("Database schema successfully applied!")
}
