package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

func GetAllPosts(db *sql.DB, tmpl *template.Template, posts []models.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fetch comments for each post
		for i, post := range posts {
			comments, err := repositories.GetComments(db, post.ID)
			if err != nil {
				log.Printf("Failed to get comments: %v", err)
				util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			posts[i].Comments = comments
		}

		// Set content type to text/html and serve the index page
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		err := tmpl.ExecuteTemplate(w, "index.html", struct {
			Posts []models.Post
		}{Posts: posts})
		if err != nil {
			log.Printf("Failed to render template: %v", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func GetAllPostsAPI(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := repositories.GetPosts(db)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// fetch comments for each post
		for i, post := range posts {
			comments, err := repositories.GetComments(db, post.ID)
			if err != nil {
				log.Printf("Failed to get posts: %v", err)
				util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			posts[i].Comments = comments
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(posts); err != nil {
			log.Printf("Failed to encode posts to JSON: %v", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// FilterPosts - Handles filtering posts by category or user
func FilterPosts(w http.ResponseWriter, r *http.Request) {
	session := StoreSession{}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Cookie not found: %v", err)
	}

	for _, v := range Sessions {
		if v.Token == cookie.Value {
			session = v
			break
		}
	}
	if r.URL.Path != "/filter" {
		log.Println("Path not found", r.URL.Path)
		util.ErrorHandler(w, "Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		log.Println("Method not allowed", r.Method)
		util.ErrorHandler(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Println("Error parsing form", err)
		util.ErrorHandler(w, "Unknown Error", http.StatusInternalServerError)
		return
	}

	categories := r.Form["category"]
	filter := r.FormValue("filter")

	if len(categories) != 0 {
		logged := true
		if session.UserId == 0 {
			logged = false
		}
		posts, err := repositories.FilterPostsByCategories(util.DB, categories)
		if err != nil {
			log.Println(err)
			util.ErrorHandler(w, "Unkown error Occured", http.StatusInternalServerError)
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
			IsLoggedIn: logged,
			Name:       "",
			Email:      "",
			Posts:      posts,
		}
		tmpl, err := template.ParseFiles("frontend/templates/index.html")
		if err != nil {
			log.Printf("Failed to load index template: %v", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		return
	}

	if session.UserId == 0 {
		log.Println("Unauthorized session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	posts := []models.Post{}

	if filter == "created" {
		posts, err = repositories.FilterPostsByUser(util.DB, session.UserId)
	}
	if filter == "liked" {
		posts, err = repositories.FilterPostsByLikes(util.DB, session.UserId)
	}
	if err != nil {
		log.Println(err)
		util.ErrorHandler(w, "Unkown error Occured", http.StatusInternalServerError)
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
		IsLoggedIn: true,
		Name:       "",
		Email:      "",
		Posts:      posts,
	}
	tmpl, err := template.ParseFiles("frontend/templates/index.html")
	if err != nil {
		log.Printf("Failed to load index template: %v", err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
