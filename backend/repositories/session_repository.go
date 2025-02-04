package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jesee-kuya/forum/backend/util"
)

// StoreSession creates a new session for a user with expiration time
func StoreSession(userID int, sessionToken string) error {
	expiration := time.Now().Add(24 * time.Hour) // Session expires in 24 hours
	query := "INSERT INTO tblSessions (user_id, session_token, expires_at) VALUES (?, ?, ?)"

	_, err := util.DB.Exec(query, userID, sessionToken, expiration)
	if err != nil {
		return fmt.Errorf("failed to store session: %v", err)
	}

	return nil
}

// ValidateSession checks if a session token is valid and not expired
func ValidateSession(sessionToken string) (int, error) {
	query := "SELECT user_id, expires_at FROM tblSessions WHERE session_token = ?"
	row := util.DB.QueryRow(query, sessionToken)

	var userID int
	var expiresAt time.Time

	err := row.Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("invalid or expired session token")
		}
		return 0, fmt.Errorf("error validating session: %v", err)
	}

	if expiresAt.Before(time.Now()) {
		_, _ = util.DB.Exec("DELETE FROM tblSessions WHERE session_token = ?", sessionToken)
		return 0, fmt.Errorf("session expired")
	}

	return userID, nil
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
