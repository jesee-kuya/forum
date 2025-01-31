package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/route"
	"github.com/jesee-kuya/forum/backend/util"
)

func main() {
	util.Init()
	defer util.DB.Close()
	
	port, err := util.ValidatePort()
	if err != nil {
		log.Fatalf("Error validating port: %v", err)
		return
	}
	router := route.InitRoutes()

	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server started at http://localhost%s\n", port)
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func getReactions() {
	reactions, err := repositories.GetReactions(util.DB, 4, "Dislike")
	if err != nil {
		fmt.Println("Could not fetch Reactions", err)
		return
	}

	fmt.Println(reactions)
}

func getFiles() {
	files, err := repositories.GetMediaFiles(util.DB, 4)
	if err != nil {
		fmt.Println("Could not fetch files", err)
		return
	}

	fmt.Println(files)
}

func addReactions() {
	reaction := models.Reaction{
		Reaction: "Dislike",
		UserID:   4,
		PostID:   4,
	}

	repositories.InsertRecord(util.DB, "tblReactions", []string{"reaction", "user_id", "post_id"}, reaction.Reaction, reaction.UserID, reaction.PostID)
}

func addPost() {
	post := models.Post{
		PostTitle:    "Monthly Goals",
		Body:         "The bocal team has brought a requirement of submitting monthly goals for every apprentice enrolled in the boot camp",
		PostCategory: "Slavery",
		UserID:       1,
	}

	repositories.InsertRecord(util.DB, "tblPosts", []string{"post_title", "body", "post_category", "user_id"}, post.PostTitle, post.Body, post.PostCategory, post.UserID)
}
