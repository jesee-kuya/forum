package handler

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

type StoreSession struct {
	Token, Email string
	UserId       int
	ExpiryTime   time.Time
}

type RequestData struct {
	ID string `json:"id"`
}

type Response struct {
	Success bool `json:"success"`
}

var (
	Session  StoreSession
	Sessions []StoreSession
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		util.ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := r.Cookie("session_token")
	if err == nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	// Load posts
	posts, err := repositories.GetPosts(util.DB)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	posts, err = PostDetails(posts)
	if err != nil {
		log.Println(err)
		util.ErrorHandler(w, "Unkown error Occured", http.StatusInternalServerError)
		return
	}
	data := struct {
		IsLoggedIn  bool
		Name, Email string
		Posts       []models.Post
	}{
		IsLoggedIn: false,
		Name:       "",
		Email:      "",
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
