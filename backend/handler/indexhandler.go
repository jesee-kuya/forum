package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

type RequestData struct {
	ID string `json:"id"`
}

type Response struct {
	Success bool `json:"success"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		util.ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	session, cookie, err := ValidateCookie(r)
	if err != nil {
		log.Printf("Failed to validate cookie: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
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

	fmt.Println("Now", time.Now())
	fmt.Println("Expiry time", session.ExpiryTime)

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
	_, err = repositories.GetUserByEmail(session.Email)
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
	PostDetails(posts, w, true, session)
}
