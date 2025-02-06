package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func PasswordEncrypt(password []byte, cost int) ([]byte, error) {
	if len(password) < 8 || len(password) > 64 {
		return nil, errors.New("password length must be between 8 and 64 characters")
	}

	// Password encryption logic
	hashedPassword, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}
