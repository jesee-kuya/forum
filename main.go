package main

import (
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/database"
	"github.com/jesee-kuya/forum/backend/handler"
	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/util"
)

func main() {
	util.Init()

	// serve static files
	fs := http.FileServer(http.Dir("./frontend/static"))
	http.Handle("/frontend/static/", http.StripPrefix("/frontend/static/", fs))

	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/sign-in", handler.LoginHandler)
	http.HandleFunc("/sign-up", handler.SignupHandler)
	http.HandleFunc("/upload", handler.UploadMedia)

	port := ":8080"
	log.Printf("Server started at http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// func addFiles() {
// 	db := database.CreateConnection()

// 	file := models.MediaFile {
// 		PostID: 1,
// 		FileName: "profile.jpg",
// 		FileType: "Profile Picture",
// 	}
// }
