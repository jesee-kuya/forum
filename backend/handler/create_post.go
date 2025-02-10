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
	var url string
	if r.Method != http.MethodPost {
		log.Println("Invalid request method:", r.Method)
		util.ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Create the img directory if it does not exist
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		log.Println("Failed to create uploads directory:", err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}

	// Parse the multipart form with a 25MB limit
	err := r.ParseMultipartForm(25 << 20)
	if err != nil {
		log.Println("Failed parsing multipart form:", err)
		util.ErrorHandler(w, "Failed parsing form data", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("uploaded-file")
	if err != nil {
		if err.Error() == "http: no such file" {
			log.Println("No file uploaded, continuing process.")
		} else {
			log.Println("Failed retrieving media file:", err)
			util.ErrorHandler(w, "File upload error", http.StatusBadRequest)
			return
		}
	}

	if file != nil {
		defer file.Close()

		log.Printf("Success in uploading %q, content type %v.\n", handler.Filename, handler.Header)

		// Validate MIME type and get the file extension
		fileExt, err := ValidateMimeType(file)
		if err != nil {
			log.Println("Invalid extension associated with file:", err)
			util.ErrorHandler(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			log.Println("Failed to reset file pointer:", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}

		// Create a temporary file with the correct extension
		tempFile, err := os.CreateTemp("uploads", "upload-*"+fileExt)
		if err != nil {
			log.Println("Failed to read file:", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			log.Println("Failed to read file:", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}

		// Write the uploaded file content to the temp file
		_, err = tempFile.Write(fileBytes)
		if err != nil {
			log.Println("Failed to write file:", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}
		tempFilePath := tempFile.Name()
		url = fmt.Sprintf("%v", tempFilePath)
	}

	cookie, err := getSessionID(r)
	if err != nil {
		log.Println("Invalid Session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sessionData, err := getSessionData(cookie)
	if err != nil {
		log.Println("Invalid Session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, err := repositories.InsertRecord(util.DB, "tblPosts", []string{"post_title", "body", "media_url", "user_id"}, r.FormValue("post-title"), r.FormValue("post-content"), url, sessionData["userId"].(int))
	if err != nil {
		log.Println("failed to add post", err)
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Println("error parsing form:", err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}

	categories := r.Form["category[]"]

	for _, category := range categories {
		repositories.InsertRecord(util.DB, "tblPostCategories", []string{"post_id", "category"}, id, category)
	}

	r.Method = http.MethodGet
	http.Redirect(w, r, "/home", http.StatusSeeOther)
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
