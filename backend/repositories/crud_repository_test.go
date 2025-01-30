package repositories

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

// setupTestDB initializes a temporary SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	tempDBFile := "test_db.db"
	db, err := sql.Open("sqlite3", tempDBFile)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create a test table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tblPosts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			content TEXT,
			status TEXT
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove(tempDBFile)
	})

	return db
}

// TestInsertRecord tests the InsertRecord function
func TestInsertRecord(t *testing.T) {
	db := setupTestDB(t)

	table := "tblPosts"
	columns := []string{"title", "content", "status"}
	values := []interface{}{"Test Title", "Test Content", "Active"}

	// Insert a record
	id, err := InsertRecord(db, table, columns, values...)
	if err != nil {
		t.Fatalf("InsertRecord failed: %v", err)
	}

	// Verify the inserted record
	var title, content, status string
	err = db.QueryRow("SELECT title, content, status FROM tblPosts WHERE id = ?", id).Scan(&title, &content, &status)
	if err != nil {
		t.Fatalf("Failed to query inserted record: %v", err)
	}

	if title != "Test Title" || content != "Test Content" || status != "Active" {
		t.Fatalf("Inserted record does not match expected values. Got: %s, %s, %s", title, content, status)
	}

	t.Logf("Successfully inserted record with ID: %d", id)
}

// TestDeleteRecord tests the DeleteRecord function
func TestDeleteRecord(t *testing.T) {
	db := setupTestDB(t)

	// Insert a test record
	table := "tblPosts"
	columns := []string{"title", "content", "status"}
	values := []interface{}{"Test Title", "Test Content", "Active"}
	id, err := InsertRecord(db, table, columns, values...)
	if err != nil {
		t.Fatalf("Failed to insert test record: %v", err)
	}

	// Delete the record
	err = DeleteRecord(db, table, "status", int(id))
	if err != nil {
		t.Fatalf("DeleteRecord failed: %v", err)
	}

	// Verify the record was marked as deleted
	var status string
	err = db.QueryRow("SELECT status FROM tblPosts WHERE id = ?", id).Scan(&status)
	if err != nil {
		t.Fatalf("Failed to query deleted record: %v", err)
	}

	if status != "Deleted" {
		t.Fatalf("Expected status 'Deleted', got: %s", status)
	}

	t.Logf("Successfully marked record with ID %d as deleted", id)
}
