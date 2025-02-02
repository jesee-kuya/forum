package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/utils"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	postType := r.URL.Query().Get("post_type")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			utils.ErrorHandler(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		limit = parsedLimit
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			utils.ErrorHandler(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
		offset = parsedOffset
	}

	posts, err := repositories.GetPostsWithPagination(postType, limit, offset)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		utils.ErrorHandler(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Printf("Error encoding posts to JSON: %v", err)
		utils.ErrorHandler(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetPostByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.ErrorHandler(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	post, err := repositories.GetPostByID(id)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		utils.ErrorHandler(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		log.Printf("Error encoding post to JSON: %v", err)
		utils.ErrorHandler(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var newPost models.Post
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		utils.ErrorHandler(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("user: %v", newPost.UserID)

	if newPost.UserID == 0 || newPost.Body == "" {
		utils.ErrorHandler(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	postID, err := repositories.CreatePost(newPost)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		utils.ErrorHandler(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"post_id": postID})
}

func GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		utils.ErrorHandler(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(userID)

	posts, err := repositories.GetPostsByUser(id)
	if err != nil {
		log.Printf("Error fetching posts for user %s: %v", userID, err)
		utils.ErrorHandler(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Printf("Error encoding posts to JSON: %v", err)
		utils.ErrorHandler(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		utils.ErrorHandler(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := 1 // Get from session/authentication (placeholder)
	err = repositories.LikePost(postID, userID)
	if err != nil {
		log.Printf("Failed to like post: %v", err)
		utils.ErrorHandler(w, "Could not like post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Post liked successfully"}`))
}

// CreateComment - Handles adding a comment to a post
func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := 1 // Get from session/authentication (placeholder)
	parentID, err := strconv.Atoi(r.URL.Query().Get("parent_id"))
	if err != nil {
		utils.ErrorHandler(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorHandler(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = repositories.CreateComment(userID, req.Body, parentID)
	if err != nil {
		log.Printf("Failed to create comment: %v", err)
		utils.ErrorHandler(w, "Could not create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Comment added successfully"}`))
}

// FilterPosts - Handles filtering posts by category or user
func FilterPosts(w http.ResponseWriter, r *http.Request) {
	filterType := r.URL.Query().Get("type")
	filterValue := r.URL.Query().Get("value")

	posts, err := repositories.FilterPosts(filterType, filterValue)
	if err != nil {
		log.Printf("Failed to filter posts: %v", err)
		utils.ErrorHandler(w, "Could not filter posts", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Printf("Failed to encode posts: %v", err)
		utils.ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
