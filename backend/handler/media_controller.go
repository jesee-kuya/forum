package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/utils"
)

// UploadMedia handles media file uploads with authentication.
func UploadMedia(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		uploaderID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
		if err != nil {
			utils.ErrorHandler(w, "Unauthorized: Missing or invalid user ID", http.StatusUnauthorized)
			return
		}

		err = r.ParseMultipartForm(50 << 20)
		if err != nil {
			utils.ErrorHandler(w, "Failed parsing form data", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			utils.ErrorHandler(w, "File upload error", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filePath, err := repositories.SaveMedia(db, file, handler.Filename, uploaderID)
		if err != nil {
			utils.ErrorHandler(w, err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		fmt.Fprintf(w, "File uploaded successfully: %s\n", filePath)
	}
}
