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

func CheckReactions(db *sql.DB, userId, postId int) (bool, string) {
	query := `
    SELECT * FROM tblReactions
    WHERE post_id = ? AND user_id = ?
`

	rows, err := db.Query(query, postId, userId)
	if err != nil {
		log.Println("error executing query:", err)
		return false, ""
	}
	defer rows.Close()

	if !rows.Next() {
		log.Println("No matching rows found")
		return false, ""
	}

	reaction := models.Reaction{}

	err = rows.Scan(&reaction.ID, &reaction.Reaction, &reaction.ReactionStatus, &reaction.UserID, &reaction.PostID)
	if err != nil {
		return false, ""
	}

	return true, reaction.Reaction
}

func UpdateReaction(db *sql.DB, reaction string, userId, postId int) {
	query := `
	UPDATE tblReactions
	SET reaction = ?
	 WHERE post_id = ? AND user_id = ?
	`
	_, err := db.Query(query, reaction, postId, userId)
	if err != nil {
		log.Println("error executing query:", err)
		return
	}
}

func UpdateReactionStatus(db *sql.DB, userId, postId int) {
	query := `
    UPDATE tblReactions
    SET reaction_status = 
        CASE 
            WHEN reaction_status = 'clicked' THEN 'unclicked'
            ELSE 'clicked'
        END
    WHERE post_id = ? AND user_id = ?
`

	_, err := db.Exec(query, postId, userId) // Use Exec instead of Query for UPDATE
	if err != nil {
		log.Println("Error executing query:", err)
		return
	}
}
