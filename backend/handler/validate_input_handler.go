package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jesee-kuya/forum/backend/util"
)

/*
ValidateInputHandler checks if a name or email already exists in the database.
*/
func ValidateInputHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/validate" {
		log.Printf("Failed finding %q route.\n", r.URL.Path)
		util.ErrorHandler(w, "Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		log.Printf("Invalid request method: %v\n", r.Method)
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Failed parsing form: %v\n", err)
		util.ErrorHandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username, email := strings.TrimSpace(r.FormValue("username")), strings.TrimSpace(r.FormValue("email"))
	var query, value string

	if username != "" {
		query = "SELECT id FROM tblUsers WHERE username = ?"
		value = username
	} else if email != "" {
		query = "SELECT id FROM tblUsers WHERE email = ?"
		value = email
	} else {
		log.Println("Invalid input provided.")
		util.ErrorHandler(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var userID int
	err := util.DB.QueryRow(query, value).Scan(&userID)
	// Provided credentials are unique
	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(map[string]bool{"available": true})
		return
	} else if err != nil {
		log.Printf("Failed quering databse: %v\n", err)
		util.ErrorHandler(w, "Something Unexpected Happened. Try Again Later", http.StatusInternalServerError)
		return
	}
	// Provided credentials are taken
	json.NewEncoder(w).Encode(map[string]bool{"available": false})
}
