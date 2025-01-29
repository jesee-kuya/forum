package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorHandler(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		utils.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = repositories.RegisterUser(req.Username, req.Email, string(hashedPassword))
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		utils.ErrorHandler(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "User registered successfully"}`))
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorHandler(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := repositories.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Failed to find user: %v", err)
		utils.ErrorHandler(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		utils.ErrorHandler(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a new UUID session token
	sessionToken, err := uuid.NewV4()
	if err != nil {
		log.Printf("Failed to generate session token: %v", err)
		utils.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store session token in database
	err = repositories.StoreSession(user.ID, sessionToken.String())
	if err != nil {
		log.Printf("Failed to store session token: %v", err)
		utils.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the session token as an HTTP cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Login successful"}`))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		utils.ErrorHandler(w, "No active session", http.StatusUnauthorized)
		return
	}

	err = repositories.DeleteSession(cookie.Value)
	if err != nil {
		utils.ErrorHandler(w, "Failed to log out", http.StatusInternalServerError)
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
