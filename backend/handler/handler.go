package handler

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
	"golang.org/x/crypto/bcrypt"
)

type StoreSession struct {
	Token, Email string
	UserId       int
	ExpiryTime   time.Time
}

var Session StoreSession

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		util.ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Cookie not found: %v", err)
		util.ErrorHandler(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	// Fetch session from DB
	dbSessionToken, err := repositories.GetSessionByUserEmail(Session.UserId)
	if err != nil || dbSessionToken != cookie.Value {
		log.Printf("Invalid session token")
		util.ErrorHandler(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	// Validate the cookie value against the session token
	if cookie.Value != Session.Token {
		log.Printf("Invalid session token: %v", err)
		util.ErrorHandler(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	} else {
		r.Method = http.MethodGet
	}

	if time.Now().After(Session.ExpiryTime) {
		log.Println("User session has expired. Please log in again")
		util.ErrorHandler(w, "User session has expired. Please log in again", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch user information
	user, err := repositories.GetUserByEmail(Session.Email)
	if err != nil {
		log.Printf("Invalid session token: %v", err)
		util.ErrorHandler(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	// Load posts
	posts, err := repositories.GetPosts(util.DB)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch comments, categories, likes, and dislikes for each post
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
		IsLoggedIn  bool
		Name, Email string
		Posts       []models.Post
	}{
		IsLoggedIn: true,
		Name:       user.Username,
		Email:      user.Email,
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sign-in" {
		util.ErrorHandler(w, "Page not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		user, err := repositories.GetUserByEmail(email)
		if err != nil {
			util.ErrorHandler(w, "Error fetching user", http.StatusForbidden)
			log.Println("Error fetching user", err)
			return
		}

		// decrypt password & authorize user
		storedPassword := user.Password

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(r.FormValue("password")))
		if err != nil {
			log.Printf("Failed to hash: %v", err)
			util.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		sessionToken := uuid.New().String()
		expiryTime := time.Now().Add(40 * time.Minute)

		err = repositories.StoreSession(user.ID, sessionToken, expiryTime)
		if err != nil {
			log.Printf("Failed to store session token: %v", err)
			util.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		Session.Token = sessionToken
		Session.UserId = user.ID
		Session.Email = user.Email
		Session.ExpiryTime = expiryTime

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  expiryTime,
			HttpOnly: true,
			Path:     "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)

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
	var user models.User
	if r.URL.Path != "/sign-up" {
		util.ErrorHandler(w, "Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		fmt.Println("OK: ", http.StatusOK)
		r.ParseForm()
		user.Username = r.PostFormValue("username")
		user.Email = r.PostFormValue("email")
		user.Password = r.PostFormValue("password")

		if user.Email == "" || user.Password == "" || user.Username == "" {
			util.ErrorHandler(w, "Fields cannot be empty", http.StatusBadRequest)
			return
		}

		hashed, err := util.PasswordEncrypt([]byte(user.Password), 10)
		if err != nil {
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		}

		id, err := repositories.InsertRecord(util.DB, "tblUsers", []string{"username", "email", "user_password"}, user.Username, user.Email, string(hashed))
		if err != nil {
			util.ErrorHandler(w, "user Can not be added", http.StatusForbidden)
			log.Println("Error adding user:", err)
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
			return
		}

		tmpl.Execute(w, nil)
	} else {
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorHandler(w, "No active session", http.StatusUnauthorized)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
