package repositories

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetCategories_QueryError(t *testing.T) {
	// Create a mock database and sqlmock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Define the expected query and simulate an error
	postID := 1
	query := `
		SELECT \* FROM tblPostCategories
		WHERE post_id = \?
	`
	mock.ExpectQuery(query).WithArgs(postID).WillReturnError(fmt.Errorf("mock query error"))

	// Call the function under test
	categories, err := GetCategories(db, postID)

	// Assert that an error is returned
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute query")

	// Assert that the result is nil
	assert.Nil(t, categories)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCategories_ScanError(t *testing.T) {
	// Create a mock database and sqlmock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Define the expected query and rows with invalid data
	postID := 1
	query := `
		SELECT \* FROM tblPostCategories
		WHERE post_id = \?
	`
	rows := sqlmock.NewRows([]string{"id", "post_id", "category_name"}).
		AddRow(1, postID, "Technology").
		AddRow("invalid", postID, "Programming") // Invalid data for ID

	// Expect the query and return the mock rows
	mock.ExpectQuery(query).WithArgs(postID).WillReturnRows(rows)

	// Call the function under test
	categories, err := GetCategories(db, postID)

	// Assert that an error is returned
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to scan row")

	// Assert that the result is nil
	assert.Nil(t, categories)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
