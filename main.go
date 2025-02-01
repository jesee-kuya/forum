package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jesee-kuya/forum/backend/database"
	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/route"
	"github.com/jesee-kuya/forum/backend/util"
)

func main() {
	util.Init()
	defer util.DB.Close()

	addPost()

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

/*
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
*/

func addPost() {
	db := database.CreateConnection()

	post := models.Post{
		ID:           1,
		UserID:       1,
		UserName:     "johnodhiambo0",
		PostTitle:    "Football Rivalry",
		Body:         "The football rivalry between Manchester United and Arsenal does not seem to end anytime soon 🔥.",
		PostCategory: "Sports",
		Likes:        15,
		Dislikes:     5,
		CommentCount: 5,
	}

	repositories.InsertRecord(db, "tblPosts", []string{"id", "user_id", "username", "post_title", "body", "post_category", "likes", "dislikes", "comment_count"}, post.ID, post.UserID, post.UserName, post.PostTitle, post.Body, post.PostCategory, post.Likes, post.Dislikes, post.CommentCount)
}
