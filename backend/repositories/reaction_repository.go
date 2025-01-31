package repositories

import (
	"database/sql"
	"fmt"

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
