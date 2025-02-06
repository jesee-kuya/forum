package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

type RequestData struct {
	ID string `json:"id"`
}

var (
	Session  StoreSession
	Sessions []StoreSession
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		util.ErrorHandler(w, "Page not found", http.StatusNotFound)
		return
	}

	posts, err := loadPosts()
	if err != nil {
		log.Println("Error loading posts:", err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsLoggedIn bool
		Name       string
		Email      string
		Posts      []models.Post
	}{false, "", "", posts}

	renderTemplate(w, "frontend/templates/index.html", data)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var session StoreSession

	if r.URL.Path != "/home" {
		util.ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Cookie not found: %v", err)
		util.ErrorHandler(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	for _, v := range Sessions {
		if v.Token == cookie.Value {
			session = v
			break
		}
	}

	// Fetch session from DB
	dbSessionToken, err := repositories.GetSessionByUserEmail(session.UserId)
	if err != nil || dbSessionToken != cookie.Value {
		log.Printf("Invalid session token: %v\n", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Validate the cookie value against the session token
	if cookie.Value != session.Token {
		log.Printf("Invalid session token: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		r.Method = http.MethodGet
	}

	if time.Now().After(session.ExpiryTime) {
		log.Println("User session has expired. Please log in again")
		util.ErrorHandler(w, "User session has expired. Please log in again", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch user information
	user, err := repositories.GetUserByEmail(session.Email)
	if err != nil {
		log.Printf("Invalid session token: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		if err1 != nil {
			log.Println("Failed to get comments:", err1)
			return
		}
		categories, err3 := repositories.GetCategories(util.DB, post.ID)
		if err3 != nil {
			log.Println("Failed to get categories", err3)
			util.ErrorHandler(w, "An unexpected error occured", http.StatusInternalServerError)
			return
		}
		likes, err4 := repositories.GetReactions(util.DB, post.ID, "Like")
		if err4 != nil {
			log.Println("Failed to get likes", err4)
			util.ErrorHandler(w, "An unexpected error occured", http.StatusInternalServerError)
			return
		}
		dislikes, err := repositories.GetReactions(util.DB, post.ID, "Dislike")
		if err != nil {
			log.Printf("Failed to get dislikes: %v", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		posts[i].Comments = comments
		posts[i].CommentCount = len(comments)
		posts[i].Categories = categories
		posts[i].Likes = len(likes)
		posts[i].Dislikes = len(dislikes)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiryTime,
		HttpOnly: true,
		Path:     "/upload",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiryTime,
		HttpOnly: true,
		Path:     "/logout",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiryTime,
		HttpOnly: true,
		Path:     "/comments",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiryTime,
		HttpOnly: true,
		Path:     "/like",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiryTime,
		HttpOnly: true,
		Path:     "/dislike",
	})

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

type Response struct {
	Success bool `json:"success"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if r.URL.Path != "/sign-up" {
		util.ErrorHandler(w, "Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		fmt.Println("OK: ", http.StatusOK)

		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed parsing form: %v\n", err)
			util.ErrorHandler(w, "Failed parsing form", http.StatusInternalServerError)
			return
		}

		user.Username = strings.TrimSpace(r.PostFormValue("username"))
		user.Email = strings.TrimSpace(r.PostFormValue("email"))
		user.Password = strings.TrimSpace(r.PostFormValue("password"))

		err = util.ValidateFormFields(user.Username, user.Email, user.Password)
		if err != nil {
			log.Printf("Invalid form values from user: %v\n", err)
			response := Response{Success: false}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		hashed, err := util.PasswordEncrypt([]byte(user.Password), 10)
		if err != nil {
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Failed encrypting password: %v\n", err)
			return
		}

		_, err = repositories.InsertRecord(util.DB, "tblUsers", []string{"username", "email", "user_password"}, user.Username, user.Email, string(hashed))
		if err != nil {
			util.ErrorHandler(w, "user cannot be added", http.StatusForbidden)
			log.Println("Error adding user:", err)
			return
		}
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sign-in" {
		util.ErrorHandler(w, "Page not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodPost:
		email := r.FormValue("email")
		user, err := repositories.GetUserByEmail(email)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.FormValue("password"))) != nil {
			util.ErrorHandler(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		sessionToken := uuid.New().String()
		expiryTime := time.Now().Add(24 * time.Hour)

		repositories.DeleteSessionByUser(user.ID)
		if err := repositories.StoreSession(user.ID, sessionToken, expiryTime); err != nil {
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		Sessions = append(Sessions, StoreSession{sessionToken, user.Email, user.ID, expiryTime})

		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: sessionToken, Expires: expiryTime, HttpOnly: true, Path: "/home"})
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	case http.MethodGet:
		renderTemplate(w, "frontend/templates/sign-in.html", nil)
	default:
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
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

	for i, session := range Sessions {
		if session.Token == cookie.Value {
			Sessions[i] = StoreSession{}
			break
		}
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

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comments" {
		log.Println("url not found", r.URL.Path)
		util.ErrorHandler(w, "Not Found", http.StatusNotFound)
		return
	}
	var session StoreSession
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Cookie not found: %v", err)
		util.ErrorHandler(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	for _, v := range Sessions {
		if v.Token == cookie.Value {
			session = v
			break
		}
	}

	if r.Method != http.MethodPost {
		log.Println("Method not allowed in comment handler", r.Method)
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.FormValue("id")
	userId := session.UserId
	comment := r.FormValue("comment")

	repositories.InsertRecord(util.DB, "tblPosts", []string{"user_id", "body", "parent_id", "post_title"}, userId, comment, id, "comment")
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/like" {
		log.Println("url not found", r.URL.Path)
		util.ErrorHandler(w, "Not Found", http.StatusNotFound)
		return
	}

	var session StoreSession
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Cookie not found: %v", err)
		util.ErrorHandler(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	for _, v := range Sessions {
		if v.Token == cookie.Value {
			session = v
			break
		}
	}

	if r.Method != http.MethodPost {
		log.Println("Method not allowed in reactions", r.Method)
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData RequestData
	err = json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		log.Println("Failed to decode json:", err)
		util.ErrorHandler(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(reqData.ID)
	if err != nil {
		log.Println("Failed to change to int:", postID)
		util.ErrorHandler(w, "An unexpected error occurred", http.StatusInternalServerError)
		return
	}

	status := "like"

	check, reaction := repositories.CheckReactions(util.DB, session.UserId, postID)
	log.Printf("CheckReactions: check=%v, reaction=%s", check, reaction) // Debugging

	if !check {
		log.Println("Inserting new reaction record") // Debugging
		_, err := repositories.InsertRecord(util.DB, "tblReactions", []string{"user_id", "post_id", "reaction"}, session.UserId, postID, status)
		if err != nil {
			log.Println("Failed to insert record:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if status == reaction {
		log.Println("Updating reaction status") // Debugging
		err := repositories.UpdateReactionStatus(util.DB, session.UserId, postID)
		if err != nil {
			log.Println("Failed to update reaction status:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	} else {
		log.Println("Updating reaction") // Debugging
		err := repositories.UpdateReaction(util.DB, status, session.UserId, postID)
		if err != nil {
			log.Println("Failed to update reaction:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dislike" {
		log.Println("url not found", r.URL.Path)
		util.ErrorHandler(w, "Not Found", http.StatusNotFound)
		return
	}

	var session StoreSession
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Cookie not found: %v", err)
		util.ErrorHandler(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	for _, v := range Sessions {
		if v.Token == cookie.Value {
			session = v
			break
		}
	}

	if r.Method != http.MethodPost {
		log.Println("Method not allowed in reactions", r.Method)
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData RequestData
	err = json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		log.Println("Failed to decode json:", err)
		util.ErrorHandler(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(reqData.ID)
	if err != nil {
		log.Println("Failed to change to int:", postID)
		util.ErrorHandler(w, "An unexpected error occurred", http.StatusInternalServerError)
		return
	}

	status := "dislike"

	check, reaction := repositories.CheckReactions(util.DB, session.UserId, postID)
	log.Printf("CheckReactions: check=%v, reaction=%s", check, reaction) // Debugging

	if !check {
		log.Println("Inserting new reaction record") // Debugging
		_, err := repositories.InsertRecord(util.DB, "tblReactions", []string{"user_id", "post_id", "reaction"}, session.UserId, postID, status)
		if err != nil {
			log.Println("Failed to insert record:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if status == reaction {
		log.Println("Updating reaction status")
		err := repositories.UpdateReactionStatus(util.DB, session.UserId, postID)
		if err != nil {
			log.Println("Failed to update reaction status:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	} else {
		log.Println("Updating reaction")
		err := repositories.UpdateReaction(util.DB, status, session.UserId, postID)
		if err != nil {
			log.Println("Failed to update reaction:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
}

func loadPosts() ([]models.Post, error) {
	posts, err := repositories.GetPosts(util.DB)
	if err != nil {
		return nil, err
	}

	for i, post := range posts {
		comments, err1 := repositories.GetComments(util.DB, post.ID)
		categories, err2 := repositories.GetCategories(util.DB, post.ID)
		likes, err3 := repositories.GetReactions(util.DB, post.ID, "Like")
		dislikes, err4 := repositories.GetReactions(util.DB, post.ID, "Dislike")
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return nil, fmt.Errorf("failed to fetch post details")
		}

		posts[i].Comments = comments
		posts[i].CommentCount = len(comments)
		posts[i].Categories = categories
		posts[i].Likes = len(likes)
		posts[i].Dislikes = len(dislikes)
	}
	return posts, nil
}

func renderTemplate(w http.ResponseWriter, templateFile string, data interface{}) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Printf("Failed to load template %s: %v", templateFile, err)
		util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
