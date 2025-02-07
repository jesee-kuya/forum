package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
	"golang.org/x/crypto/bcrypt"
)

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

		sessionToken, err := uuid.NewV4()
		if err != nil {
			log.Printf("Failed to get uuid: %v", err)
			util.ErrorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		expiryTime := time.Now().UTC().Add(24 * time.Hour)

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

		CookieSession(w, Session)
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
