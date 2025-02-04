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


func addCategories() {
	cat := models.Category{
		PostID: 7,
		CategoryName: "Technology",
	}

	repositories.InsertRecord(util.DB, "tblPostCategories", []string{"post_id", "category"}, cat.PostID, cat.CategoryName)
}

func addPost() {
	post := models.Post{
		PostTitle:    "Sunday Rost",
		Body:         "The earth will be exactly six thousand years in 2027",
		UserID:       1,
	}

	id, err := repositories.InsertRecord(util.DB, "tblPosts", []string{"post_title", "body", "user_id"}, post.PostTitle, post.Body, post.UserID)

	if err != nil {
		fmt.Println("failed to AD post", err)
		return
	}

	fmt.Println("added post", id)
}