package repositories

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jesee-kuya/forum/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestGetReactions(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define the expected query and its result
	postID := 1
	reactionType := "like"
	expectedReactions := []models.Reaction{
		{ID: 1, Reaction: "like", ReactionStatus: "clicked", UserID: 1, PostID: 1},
		{ID: 2, Reaction: "like", ReactionStatus: "clicked", UserID: 2, PostID: 1},
	}

	rows := sqlmock.NewRows([]string{"id", "reaction", "reaction_status", "user_id", "post_id"}).
		AddRow(expectedReactions[0].ID, expectedReactions[0].Reaction, expectedReactions[0].ReactionStatus, expectedReactions[0].UserID, expectedReactions[0].PostID).
		AddRow(expectedReactions[1].ID, expectedReactions[1].Reaction, expectedReactions[1].ReactionStatus, expectedReactions[1].UserID, expectedReactions[1].PostID)

	mock.ExpectQuery("SELECT \\* FROM tblReactions WHERE post_id = \\? AND reaction = \\? AND reaction_status = 'clicked'").
		WithArgs(postID, reactionType).
		WillReturnRows(rows)

	// Call the function under test
	reactions, err := GetReactions(db, postID, reactionType)
	assert.NoError(t, err)
	assert.Equal(t, expectedReactions, reactions)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetReactions_Error(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define the expected query and simulate an error
	postID := 1
	reactionType := "like"

	mock.ExpectQuery("SELECT \\* FROM tblReactions WHERE post_id = \\? AND reaction = \\? AND reaction_status = 'clicked'").
		WithArgs(postID, reactionType).
		WillReturnError(fmt.Errorf("some error"))

	// Call the function under test
	reactions, err := GetReactions(db, postID, reactionType)
	assert.Error(t, err)
	assert.Nil(t, reactions)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
