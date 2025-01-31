package repositories

import (
	"database/sql"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jesee-kuya/forum/backend/models"
	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB initializes a temporary SQLite database for testing
func setupTestDBP(t *testing.T) *sql.DB {
	// Create a temporary database file
	tempDBFile := "test_forum.db"
	db, err := sql.Open("sqlite3", tempDBFile)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create the necessary tables and insert test data
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tblUsers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS tblPosts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_title TEXT,
			body TEXT,
			created_on DATETIME,
			post_category TEXT,
			parent_id INTEGER,
			post_status TEXT DEFAULT 'visible',
			FOREIGN KEY (user_id) REFERENCES tblUsers(id)
		);

		INSERT INTO tblUsers (username) VALUES ('user1'), ('user2');

		INSERT INTO tblPosts (user_id, post_title, body, created_on, post_category, parent_id, post_status)
		VALUES
			(1, 'Post 1', 'Content 1', '2023-10-01 10:00:00', 'General', NULL, 'visible'),
			(2, 'Post 2', 'Content 2', '2023-10-02 11:00:00', 'Tech', NULL, 'visible'),
			(1, 'Comment 1', 'Comment Content 1', '2023-10-01 10:30:00', 'General', 1, 'visible'),
			(2, 'Comment 2', 'Comment Content 2', '2023-10-02 11:30:00', 'Tech', 2, 'visible');
	`)
	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}

	// Clean up the database file after the test
	t.Cleanup(func() {
		db.Close()
		os.Remove(tempDBFile)
	})

	return db
}

// TestGetPosts tests the GetPosts function
func TestGetPosts(t *testing.T) {
	db := setupTestDBP(t)

	// Call GetPosts
	posts, err := GetPosts(db)
	if err != nil {
		t.Fatalf("GetPosts failed: %v", err)
	}

	// Verify the results
	expectedPosts := []models.Post{
		{
			ID:           1,
			UserID:       1,
			UserName:     "user1",
			PostTitle:    "Post 1",
			Body:         "Content 1",
			CreatedOn:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			PostCategory: "General",
		},
		{
			ID:           2,
			UserID:       2,
			UserName:     "user2",
			PostTitle:    "Post 2",
			Body:         "Content 2",
			CreatedOn:    time.Date(2023, 10, 2, 11, 0, 0, 0, time.UTC),
			PostCategory: "Tech",
		},
	}

	if len(posts) != len(expectedPosts) {
		t.Fatalf("Expected %d posts, got %d", len(expectedPosts), len(posts))
	}

	for i, post := range posts {
		if !reflect.DeepEqual(post, expectedPosts[i]) {
			t.Errorf("Post %d does not match expected values. Got: %+v, Expected: %+v", i+1, post, expectedPosts[i])
		}
	}
}

// TestGetComments tests the GetComments function
func TestGetComments(t *testing.T) {
	db := setupTestDBP(t)

	// Call GetComments for post ID 1
	comments, err := GetComments(db, 1)
	if err != nil {
		t.Fatalf("GetComments failed: %v", err)
	}

	// Verify the results
	expectedComments := []models.Post{
		{
			ID:           3,
			UserID:       1,
			UserName:     "user1",
			PostTitle:    "Comment 1",
			Body:         "Comment Content 1",
			CreatedOn:    time.Date(2023, 10, 1, 10, 30, 0, 0, time.UTC),
			PostCategory: "General",
		},
	}

	if len(comments) != len(expectedComments) {
		t.Fatalf("Expected %d comments, got %d", len(expectedComments), len(comments))
	}

	for i, comment := range comments {
		if !reflect.DeepEqual(comment, expectedComments[i]) {
			t.Errorf("Comment %d does not match expected values. Got: %+v, Expected: %+v", i+1, comment, expectedComments[i])
		}
	}
}

// TestProcessSQLData tests the processSQLData function
func TestProcessSQLData(t *testing.T) {
	db := setupTestDBP(t)

	// Query rows from the database
	rows, err := db.Query(`
		SELECT p.id, p.user_id, u.username, p.post_title, p.body, p.created_on, p.post_category
		FROM tblPosts p
		JOIN tblUsers u ON p.user_id = u.id
		WHERE p.parent_id IS NULL AND p.post_status = 'visible'
	`)
	if err != nil {
		t.Fatalf("Failed to query rows: %v", err)
	}
	defer rows.Close()

	// Call processSQLData
	posts, err := processSQLData(rows)
	if err != nil {
		t.Fatalf("processSQLData failed: %v", err)
	}

	// Verify the results
	expectedPosts := []models.Post{
		{
			ID:           1,
			UserID:       1,
			UserName:     "user1",
			PostTitle:    "Post 1",
			Body:         "Content 1",
			CreatedOn:    time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			PostCategory: "General",
		},
		{
			ID:           2,
			UserID:       2,
			UserName:     "user2",
			PostTitle:    "Post 2",
			Body:         "Content 2",
			CreatedOn:    time.Date(2023, 10, 2, 11, 0, 0, 0, time.UTC),
			PostCategory: "Tech",
		},
	}

	if len(posts) != len(expectedPosts) {
		t.Fatalf("Expected %d posts, got %d", len(expectedPosts), len(posts))
	}

	for i, post := range posts {
		if !reflect.DeepEqual(post, expectedPosts[i]) {
			t.Errorf("Post %d does not match expected values. Got: %+v, Expected: %+v", i+1, post, expectedPosts[i])
		}
	}
}
