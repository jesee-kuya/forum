package repositories

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jesee-kuya/forum/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestGetMediaFiles_Success(t *testing.T) {
	// Create a mock database and sqlmock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Define the expected query and rows
	postID := 1
	query := `
		SELECT \* FROM tblMediaFiles 
		WHERE post_id = \? AND file_status = 'visible'
	`
	rows := sqlmock.NewRows([]string{"id", "post_id", "file_name", "file_type", "file_status"}).
		AddRow(1, postID, "image1.jpg", "image/jpeg", "visible").
		AddRow(2, postID, "video1.mp4", "video/mp4", "visible")

	// Expect the query and return the mock rows
	mock.ExpectQuery(query).WithArgs(postID).WillReturnRows(rows)

	// Call the function under test
	files, err := GetMediaFiles(db, postID)

	// Assert that there are no errors
	assert.NoError(t, err)

	// Assert the expected results
	expectedFiles := []models.MediaFile{
		{ID: 1, PostID: postID, FileName: "image1.jpg", FileType: "image/jpeg", FileStatus: "visible"},
		{ID: 2, PostID: postID, FileName: "video1.mp4", FileType: "video/mp4", FileStatus: "visible"},
	}
	assert.Equal(t, expectedFiles, files)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMediaFiles_QueryError(t *testing.T) {
	// Create a mock database and sqlmock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Define the expected query and simulate an error
	postID := 1
	query := `
		SELECT \* FROM tblMediaFiles 
		WHERE post_id = \? AND file_status = 'visible'
	`
	mock.ExpectQuery(query).WithArgs(postID).WillReturnError(fmt.Errorf("mock query error"))

	// Call the function under test
	files, err := GetMediaFiles(db, postID)

	// Assert that an error is returned
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute query")

	// Assert that the result is nil
	assert.Nil(t, files)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMediaFiles_ScanError(t *testing.T) {
	// Create a mock database and sqlmock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Define the expected query and rows with invalid data
	postID := 1
	query := `
		SELECT \* FROM tblMediaFiles 
		WHERE post_id = \? AND file_status = 'visible'
	`
	rows := sqlmock.NewRows([]string{"id", "post_id", "file_name", "file_type", "file_status"}).
		AddRow(1, postID, "image1.jpg", "image/jpeg", "visible").
		AddRow("invalid", postID, "video1.mp4", "video/mp4", "visible") // Invalid data for ID

	// Expect the query and return the mock rows
	mock.ExpectQuery(query).WithArgs(postID).WillReturnRows(rows)

	// Call the function under test
	files, err := GetMediaFiles(db, postID)

	// Assert that an error is returned
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to scan row")

	// Assert that the result is nil
	assert.Nil(t, files)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
