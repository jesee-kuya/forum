package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jesee-kuya/forum/backend/models"
)

func GetReactions(db *sql.DB, id int, react string) ([]models.Reaction, error) {
	query := `
		SELECT * FROM tblReactions
		WHERE post_id = ? AND reaction = ? AND reaction_status = 'clicked'
	`
	rows, err := db.Query(query, id, react)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var Reactions []models.Reaction

	for rows.Next() {
		reaction := models.Reaction{}

		err := rows.Scan(&reaction.ID, &reaction.Reaction, &reaction.ReactionStatus, &reaction.UserID, &reaction.PostID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		Reactions = append(Reactions, reaction)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return Reactions, nil
}

func UpdateLikeCount(db *sql.DB, postID int, userID int) (int, error) {
	query := `
			INSERT INTO tblReactions (reaction, reaction_status, user_id, post_id)
			VALUES ('like', 'clicked', ?, ?)
			ON CONFLICT(user_id, post_id, reaction) 
			DO UPDATE SET reaction_status = 'clicked'
	`

	_, err := db.Exec(query, userID, postID)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return 0, fmt.Errorf("failed to update like count: %w", err)
	}

	// Retrieve updated like count
	var newLikeCount int
	countQuery := `SELECT COUNT(*) FROM tblReactions WHERE post_id = ? AND reaction = 'like'`
	err = db.QueryRow(countQuery, postID).Scan(&newLikeCount)
	if err != nil {
		log.Printf("Error fetching new like count: %v\n", err)
		return 0, fmt.Errorf("failed to fetch new like count: %w", err)
	}

	log.Printf("New like count: %d\n", newLikeCount)
	return newLikeCount, nil
}
