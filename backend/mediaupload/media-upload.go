package mediaupload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

/*
UploadMedia handler function is responsible for performing server operations to enable media upload with a limit of up to 25 mbs.
*/
func UploadMedia(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("frontend/templates/media-upload.html")
	if err != nil {
		http.Error(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		log.Println("Failed parsing templates:", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		log.Println("Failed executing template:", err)
		return
	}

	if r.FormValue("upload") != "" {
		fmt.Fprintf(w, "Uploading file...\n")

		// Parse the multipart form with a 25MB limit
		err := r.ParseMultipartForm(25 << 20)
		if err != nil {
			http.Error(w, "Failed parsing form data", http.StatusBadRequest)
			log.Println("Failed parsing multipart form:", err)
			return
		}

		file, handler, err := r.FormFile("post1")
		if err != nil {
			http.Error(w, "File upload error", http.StatusBadRequest)
			log.Println("Failed retrieving media file:", err)
			return
		}
		defer file.Close()

		fmt.Printf("Successfully uploaded file: %v\n", handler.Filename)

		// Create a temporary file
		tempFile, err := os.CreateTemp("img", "upload-*.jpeg")
		if err != nil {
			http.Error(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			log.Println("Failed creating a temporary file:", err)
			return
		}
		defer tempFile.Close()

		// Copy file contents to the temporary file
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			log.Println("Failed saving file to temporary location:", err)
			return
		}
		fmt.Fprintf(w, "Successfully uploaded file.\n")
	}
}
