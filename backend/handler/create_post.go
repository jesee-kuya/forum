package handler

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

/*
UploadMedia handler function is responsible for performing server operations to enable media upload with a file size limit of up to 25 mbs.
*/
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var session StoreSession
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Cookie not found: %v", err)
		util.ErrorHandler(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	var url string

	for _, v := range Sessions {
		if v.Token == cookie.Value {
			session = v
			break
		}
	}

	if r.Method != http.MethodPost {
		util.ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Println("Invalid request method:", r.Method)
		return
	}

	// Create the img directory if it does not exist
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		log.Println("Failed to create uploads directory:", err)
		return
	}

	// Parse the multipart form with a 25MB limit
	err = r.ParseMultipartForm(25 << 20)
	if err != nil {
		util.ErrorHandler(w, "Failed parsing form data", http.StatusBadRequest)
		log.Println("Failed parsing multipart form:", err)
		return
	}

	file, handler, err := r.FormFile("uploaded-file")
	if err != nil {
		if err.Error() == "http: no such file" {
			log.Println("No file uploaded, continuing process.")
		} else {
			util.ErrorHandler(w, "File upload error", http.StatusBadRequest)
			log.Println("Failed retrieving media file:", err)
			return
		}
	}

	if file != nil {
		defer file.Close()

		log.Printf("Success in uploading %q, content type %v.\n", handler.Filename, handler.Header)

		// Validate MIME type and get the file extension
		fileExt, err := ValidateMimeType(file)
		if err != nil {
			util.ErrorHandler(w, err.Error(), http.StatusBadRequest)
			log.Println("Invalid extension associated with file:", err)
			return
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			util.ErrorHandler(w, "Failed to reset file pointer", http.StatusInternalServerError)
			log.Println("Failed to reset file pointer:", err)
			return
		}

		// Create a temporary file with the correct extension
		tempFile, err := os.CreateTemp("uploads", "upload-*"+fileExt)
		if err != nil {
			util.ErrorHandler(w, "Failed to create file", http.StatusInternalServerError)
			log.Println("Failed to read file:", err)
			return
		}
		defer tempFile.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			util.ErrorHandler(w, "Failed to read file", http.StatusInternalServerError)
			log.Println("Failed to read file:", err)
			return
		}

		// Write the uploaded file content to the temp file
		_, err = tempFile.Write(fileBytes)
		if err != nil {
			util.ErrorHandler(w, "Failed to write file", http.StatusInternalServerError)
			log.Println("Failed to write file:", err)
			return
		}
		tempFilePath := tempFile.Name()
		url = fmt.Sprintf("%v", tempFilePath)
	}

	id, err := repositories.InsertRecord(util.DB, "tblPosts", []string{"post_title", "body", "media_url", "user_id"}, r.FormValue("post-title"), r.FormValue("post-content"), url, session.UserId)
	if err != nil {
		fmt.Println("failed to add post", err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	categories := r.Form["category[]"]

	for _, category := range categories {
		repositories.InsertRecord(util.DB, "tblPostCategories", []string{"post_id", "category"}, id, category)
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
	r.Method = http.MethodGet
	IndexHandler(w, r)
}

/*
ValidateMimeType is used to check the MIME type of an uploaded file. It returns the extension associated with the file.
*/
func ValidateMimeType(file multipart.File) (string, error) {
	allowedMIMEs := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
		"image/webp": ".webp",
	}

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		log.Printf("Failed to read buffer: %v\n", err)
		return "", fmt.Errorf("failed to read file data")
	}

	mimeType := http.DetectContentType(buffer)
	ext, valid := allowedMIMEs[mimeType]
	if !valid {
		return "", fmt.Errorf("invalid file type")
	}
	return ext, nil
}
