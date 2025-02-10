package handler

import (
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
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

	cookie, _ := getSessionID(r)
	_, ok := SessionStore[cookie]
	if ok {
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

	PostDetails(w, r, posts, false)
}
