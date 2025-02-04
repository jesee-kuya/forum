package handler

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

var User models.User

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		util.ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		fmt.Println("OK: ", http.StatusOK)
	} else {
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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
		comments, err1 := repositories.GetComments(util.DB, post.ID)
		categories, err3 := repositories.GetCategories(util.DB, post.ID)
		likes, err4 := repositories.GetReactions(util.DB, post.ID, "Like")
		dislikes, err := repositories.GetReactions(util.DB, post.ID, "Dislike")
		if err != nil || err1 != nil || err3 != nil || err4 != nil {
			log.Printf("Failed to get posts details: %v", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		posts[i].Comments = comments
		posts[i].CommentCount = len(comments)
		posts[i].Categories = categories
		posts[i].Likes = len(likes)
		posts[i].Dislikes = len(dislikes)

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

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		User, err := repositories.GetUserByEmail(email)
		if err != nil {
			util.ErrorHandler(w, "Error fetching User", http.StatusForbidden)
			log.Println("Error fetching User", err)
			return
		}
		fmt.Printf("User: %v", User.Email)

		sessionToken, err := uuid.NewV4()
		if err != nil {
			log.Printf("Failed to generate session token: %v", err)
			util.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		err = repositories.StoreSession(User.ID, sessionToken.String())
		if err != nil {
			log.Printf("Failed to store session token: %v", err)
			util.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken.String(),
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		r.Method = http.MethodGet
		IndexHandler(w, r)
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

	if r.Method == http.MethodPost {
		r.ParseForm()
		User.Username = r.PostFormValue("username")
		User.Email = r.PostFormValue("email")
		User.Password = r.PostFormValue("password")

		if User.Email == "" || User.Password == "" || User.Username == "" {
			util.ErrorHandler(w, "Fields cannot be empty", http.StatusBadRequest)
			log.Println("Empty field")
			return
		}
		id, err := repositories.InsertRecord(util.DB, "tblUsers", []string{"Username", "email", "User_password"}, User.Username, User.Email, User.Password)
		if err != nil {
			util.ErrorHandler(w, "User Can not be added", http.StatusForbidden)
			log.Println("Error adding User:", err)
			return
		}
		fmt.Println(id)

		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		r.Method = http.MethodGet
		SignupHandler(w, r)
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

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		util.ErrorHandler(w, "No active session", http.StatusUnauthorized)
		return
	}

	err = repositories.DeleteSession(cookie.Value)
	if err != nil {
		util.ErrorHandler(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logged out successfully"}`))
}
