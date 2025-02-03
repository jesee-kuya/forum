package handler

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("post-title")
	if title != "" {
		AddPostHandler(w, r)
	}
	// URL path
	if r.URL.Path != "/" {
		util.ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	// Method used
	if r.Method == http.MethodGet {
		fmt.Println("OK: ", http.StatusOK)
	} else {
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// template rendering
	tmpl, err := template.ParseFiles("frontend/templates/index.html")
	if err != nil {
		log.Printf("Failed to load index template: %v", err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	posts, err := repositories.GetPosts(util.DB)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// fetch comments for each post
	for i, post := range posts {
		comments, err := repositories.GetComments(util.DB, post.ID)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		posts[i].Comments = comments
	}

	data := struct {
		Posts []models.Post
	}{
		Posts: posts,
	}
	tmpl.Execute(w, data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sign-in" {
		util.ErrorHandler(w, "Page not found", http.StatusNotFound)
		return
	}

	var user models.User

	if r.Method == http.MethodPost {
		if user.Username == "" || user.Email == "" || user.Password == "" {
			util.ErrorHandler(w, "Username/Email and Password are required", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return

	} else if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("frontend/templates/sign-in.html")
		if err != nil {
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		}

		tmpl.Execute(w, nil)
	} else {
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sign-up" {
		util.ErrorHandler(w, "Page Not Found", http.StatusNotFound)
		return
	}

	var user models.User

	if r.Method == http.MethodPost {
		fmt.Println("OK: ", http.StatusOK)
		r.ParseForm()
		user.Username = r.PostFormValue("username")
		user.Email = r.PostFormValue("email")
		user.Password = r.PostFormValue("password")

		if user.Username == "" || user.Email == "" || user.Password == "" {
			util.ErrorHandler(w, "Fields cannot be empty", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	} else if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("frontend/templates/sign-up.html")
		if err != nil {
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		}

		tmpl.Execute(w, nil)
	} else {
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func AddPostHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("post-title")
	repositories.InsertRecord(util.DB, "tblPosts", []string{"post_title", "body", "post_category", "user_id"}, title, r.FormValue("post-content"), r.FormValue("category"), 1)
}
