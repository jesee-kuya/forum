package main

import (
	"fmt"

	"github.com/jesee-kuya/forum/backend/database"
	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
)

func main() {
	// util.Init()

	// // serve static files
	// fs := http.FileServer(http.Dir("./frontend/static"))
	// http.Handle("/frontend/static/", http.StripPrefix("/frontend/static/", fs))

	// http.HandleFunc("/", handler.IndexHandler)
	// http.HandleFunc("/sign-in", handler.LoginHandler)
	// http.HandleFunc("/sign-up", handler.SignupHandler)
	// http.HandleFunc("/upload", handler.UploadMedia)

	// port := ":8080"
	// log.Printf("Server started at http://localhost%s\n", port)
	// err := http.ListenAndServe(port, nil)
	// if err != nil {
	// 	log.Fatalf("Error starting server: %v", err)
	// }

	getReactions()
}

func getReactions() {
	db := database.CreateConnection()
	reactions, err := repositories.GetReactions(db, 4, "Dislike")

	if err != nil {
		fmt.Println("Could not fetch Reactions", err)
		return
	}

	fmt.Println(reactions)
}

func getFiles() {
	db := database.CreateConnection()
	files, err := repositories.GetMediaFiles(db, 4)

	if err != nil {
		fmt.Println("Could not fetch files", err)
		return
	}

	fmt.Println(files)
}

func addReactions() {
	db := database.CreateConnection()

	reaction := models.Reaction{
		Reaction: "Dislike",
		UserID:   4,
		PostID:   4,
	}

	repositories.InsertRecord(db, "tblReactions", []string{"reaction", "user_id", "post_id"}, reaction.Reaction, reaction.UserID, reaction.PostID)
}
