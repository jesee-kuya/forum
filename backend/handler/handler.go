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

func PostDetails(posts []models.Post) ([]models.Post, error) {
	for i, post := range posts {
		comments, err1 := repositories.GetComments(util.DB, post.ID)
		if err1 != nil {
			log.Println("Failed to get comments:", err1)
			return nil, fmt.Errorf("failed to get comments: %v", err1)
		}

		// Getting comment reactions
		for j, comment := range comments {
			commentLikes, errLikes := repositories.GetReactions(util.DB, comment.ID, "Like")
			if errLikes != nil {
				log.Println("Failed to get likes", errLikes)
				return nil, fmt.Errorf("failed to get likes: %v", errLikes)
			}

			commentDislikes, errDislikes := repositories.GetReactions(util.DB, comment.ID, "Dislike")
			if errDislikes != nil {
				log.Println("Failed to get dislikes", errDislikes)
				return nil, fmt.Errorf("failed to get dislikes: %v", errDislikes)
			}

			comments[j].Likes = len(commentLikes)
			comments[j].Dislikes = len(commentDislikes)
		}

		categories, err3 := repositories.GetCategories(util.DB, post.ID)
		if err3 != nil {
			log.Println("Failed to get categories", err3)
			return nil, fmt.Errorf("failed to get categories: %v", err3)
		}
		likes, err4 := repositories.GetReactions(util.DB, post.ID, "Like")
		if err4 != nil {
			log.Println("Failed to get likes", err4)
			return nil, fmt.Errorf("failed to get likes: %v", err4)
		}
		dislikes, err := repositories.GetReactions(util.DB, post.ID, "Dislike")
		if err != nil {
			log.Printf("Failed to get dislikes: %v", err)
			return nil, fmt.Errorf("failed to get dislikes: %v", err)
		}

		posts[i].Comments = comments
		posts[i].CommentCount = len(comments)
		posts[i].Categories = categories
		posts[i].Likes = len(likes)
		posts[i].Dislikes = len(dislikes)
	}
	return posts, nil
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
		Path:     "/reaction",
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
			// util.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			response := Response{Success: false}
			json.NewEncoder(w).Encode(response)
			return
		}

		sessionToken, err := uuid.NewV6()
		if err != nil {
			log.Printf("Failed to get uuid: %v", err)
			util.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		expiryTime := time.Now().Add(1440 * time.Minute)

		err = repositories.DeleteSessionByUser(user.ID)
		if err != nil {
			log.Printf("Failed to delete session token: %v", err)
			util.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		strSessionToken := sessionToken.String()
		err = repositories.StoreSession(user.ID, strSessionToken, expiryTime)
		if err != nil {
			log.Printf("Failed to store session token: %v", err)
			util.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		Session.Token = strSessionToken
		Session.UserId = user.ID
		Session.Email = user.Email
		Session.ExpiryTime = expiryTime
		Sessions = append(Sessions, Session)
		Session = StoreSession{}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    strSessionToken,
			Expires:  expiryTime,
			HttpOnly: true,
			Path:     "/home",
		})

		http.Redirect(w, r, "/home", http.StatusSeeOther)

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

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	session := StoreSession{}
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
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err = r.ParseForm()
	if err != nil {
		util.ErrorHandler(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	reactionType := r.FormValue("reaction")
	postID, _ := strconv.Atoi(r.FormValue("post_id"))

	fmt.Println("Reaction: ", reactionType)
	fmt.Println("Post ID: ", postID)

	check, reaction := repositories.CheckReactions(util.DB, session.UserId, postID)

	if !check {
		_, err := repositories.InsertRecord(util.DB, "tblReactions", []string{"user_id", "post_id", "reaction"}, session.UserId, postID, reactionType)
		if err != nil {
			log.Println("Failed to insert record:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if reactionType == reaction {
		err := repositories.UpdateReactionStatus(util.DB, session.UserId, postID)
		if err != nil {
			log.Println("Failed to update reaction status:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	} else {
		err := repositories.UpdateReaction(util.DB, reactionType, session.UserId, postID)
		if err != nil {
			log.Println("Failed to update reaction:", err)
			util.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	r.Method = http.MethodGet
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
