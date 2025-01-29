package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jesee-kuya/forum/backend/models"
)

func RegisterUser(username, email, password string) error {
	query := "INSERT INTO tblUsers (username, email, user_password) VALUES (?, ?, ?)"
	_, err := db.Exec(query, username, email, password)
	if err != nil {
		return fmt.Errorf("failed to register user: %v", err)
	}
	return nil
}

func GetUserByEmail(email string) (models.User, error) {
	query := "SELECT id, username, email, user_password FROM tblUsers WHERE email = ?"
	row := db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("failed to retrieve user: %v", err)
	}
	return user, nil
}
