package test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jesee-kuya/forum/backend/repositories"
)

// TestGetPosts tests the GetPosts function
func TestGetPosts(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer mockDB.Close()

	repositories.Db = mockDB

	t.Run("Failure - Query Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM tblposts").WillReturnError(errors.New("query failed"))

		_, err := repositories.GetPosts()
		if err == nil || err.Error() != "failed to execute query: query failed" {
			t.Errorf("Expected query error, but got: %v", err)
		}
	})
}
