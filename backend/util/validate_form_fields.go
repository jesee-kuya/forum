package util

import (
	"fmt"
	"regexp"
	"strings"
)

/*
ValidateInput checks for the validity of input values provided via the form.
*/
func ValidateFormFields(userName, email, password string) error {
	if len(strings.TrimSpace(userName)) == 0 {
		return fmt.Errorf("username field cannot be empty")
	}

	if len(strings.TrimSpace(email)) == 0 {
		return fmt.Errorf("email field cannot be empty")
	}

	if len(strings.TrimSpace(password)) == 0 {
		return fmt.Errorf("password field cannot be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email address format")
	}

	if len(strings.TrimSpace(password)) < 8 {
		return fmt.Errorf("password must contain atleast 8 characters")
	}
	return nil
}
