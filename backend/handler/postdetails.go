package handler

import (
	"log"
	"net/http"
	"text/template"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

func PostDetails(posts []models.Post, w http.ResponseWriter, logged bool, session StoreSession) {
	for i, post := range posts {
		comments, err1 := repositories.GetComments(util.DB, post.ID)
		if err1 != nil {
			log.Println("Failed to get comments:", err1)
			util.ErrorHandler(w, "Unkown error Occured", http.StatusInternalServerError)
			return
		}

		// Getting comment reactions
		for j, comment := range comments {
			commentLikes, errLikes := repositories.GetReactions(util.DB, comment.ID, "Like")
			if errLikes != nil {
				log.Println("Failed to get likes", errLikes)
				util.ErrorHandler(w, "Unkown error Occured", http.StatusInternalServerError)
				return
			}

			commentDislikes, errDislikes := repositories.GetReactions(util.DB, comment.ID, "Dislike")
			if errDislikes != nil {
				log.Println("Failed to get dislikes", errDislikes)
				util.ErrorHandler(w, "Unkown error Occured", http.StatusInternalServerError)
				return
			}

			comments[j].Likes = len(commentLikes)
			comments[j].Dislikes = len(commentDislikes)
		}

		categories, err3 := repositories.GetCategories(util.DB, post.ID)
		if err3 != nil {
			log.Println("Failed to get categories", err3)
			util.ErrorHandler(w, "Unkown error Occured", http.StatusInternalServerError)
			return
		}
		likes, err4 := repositories.GetReactions(util.DB, post.ID, "Like")
		if err4 != nil {
			log.Println("Failed to get likes", err4)
			util.ErrorHandler(w, "Unkown error Occured", http.StatusInternalServerError)
			return
		}
		dislikes, err := repositories.GetReactions(util.DB, post.ID, "Dislike")
		if err != nil {
			log.Printf("Failed to get dislikes: %v", err)
			util.ErrorHandler(w, "Unkown error Occured", http.StatusInternalServerError)
			return
		}

		posts[i].Comments = comments
		posts[i].CommentCount = len(comments)
		posts[i].Categories = categories
		posts[i].Likes = len(likes)
		posts[i].Dislikes = len(dislikes)
	}
	data := struct {
		IsLoggedIn  bool
		Name, Email string
		Posts       []models.Post
	}{
		IsLoggedIn: logged,
		Name:       "",
		Email:      session.Email,
		Posts:      posts,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("frontend/templates/index.html")
	if err != nil {
		log.Printf("Failed to load index template: %v", err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
