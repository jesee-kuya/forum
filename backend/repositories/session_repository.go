package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jesee-kuya/forum/backend/util"
)

// StoreSession creates a new session for a user with expiration time
func StoreSession(userID int, sessionToken string, expiryTime time.Time) error {
	_, err := InsertRecord(util.DB, "tblSessions", []string{"user_id", "session_token", "expires_at"}, userID, sessionToken, expiryTime)
	if err != nil {
		log.Println("Error inserting session:", err)
		return err
	}
	return nil
}

// ValidateSession checks if a session token is valid and not expired
func ValidateSession(sessionToken string) (string, error) {
	query := "SELECT user_id, expires_at FROM tblSessions WHERE session_token = ?"
	row := util.DB.QueryRow(query, sessionToken)

	var sessionTOken string
	var expiresAt time.Time

	err := row.Scan(&sessionToken, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid or expired session token")
		}
		return "", fmt.Errorf("error validating session: %v", err)
	}

	if expiresAt.Before(time.Now()) {
		_, _ = util.DB.Exec("DELETE FROM tblSessions WHERE session_token = ?", sessionToken)
		return "", fmt.Errorf("session expired")
	}
	return sessionTOken, nil
}

// DeleteSession removes a session when a user logs out
func DeleteSession(sessionToken string) error {
	query := "DELETE FROM tblSessions WHERE session_token = ?"
	_, err := util.DB.Exec(query, sessionToken)
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}
	return nil
}

func DeleteSessionByUser(userId int) error {
	query := "DELETE FROM tblSessions WHERE user_id = ?"
	_, err := util.DB.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}
	return nil
}

// GetSessionByUserEmail fetches the session token associated with a given user email.
func GetSessionByUserEmail(user_id int) (string, error) {
	var sessionToken string

	query := "SELECT session_token FROM tblSessions WHERE user_id = ?"
	err := util.DB.QueryRow(query, user_id).Scan(&sessionToken)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No session found for user:", user_id)
			return "", nil
		}
		log.Println("Error retrieving session:", err)
		return "", err
	}
	return sessionToken, nil
}
