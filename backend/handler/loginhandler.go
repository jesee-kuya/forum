package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
	"golang.org/x/crypto/bcrypt"
)

var SessionStore = make(map[string]map[string]interface{})

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var err error
	if r.URL.Path != "/sign-in" {
		util.ErrorHandler(w, "Page not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		if isValidEmail(email) {
			user, err = repositories.GetUserByEmail(email)
			if err != nil {
				util.ErrorHandler(w, "Error fetching user", http.StatusForbidden)
				log.Println("Error fetching user", err)
				return
			}
		} else {
			user, err = repositories.GetUserByName(email)
			if err != nil {
				util.ErrorHandler(w, "Error fetching user", http.StatusForbidden)
				log.Println("Error fetching user", err)
				return
			}
		}

		// decrypt password & authorize user
		storedPassword := user.Password

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(r.FormValue("password")))
		if err != nil {
			log.Printf("Failed to hash: %v", err)
			w.Header().Set("Content-Type", "application/json")
			response := Response{Success: false}
			json.NewEncoder(w).Encode(response)
			return
		}

		sessionToken := createSession()

		err = repositories.DeleteSessionByUser(user.ID)
		if err != nil {
			log.Printf("Failed to delete session token: %v", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}

		EnableCors(w)

		setSessionCookie(w, sessionToken)
		setSessionData(sessionToken, "userId", user.ID)
		setSessionData(sessionToken, "userEmail", user.Email)
		expiryTime := time.Now().Add(24 * time.Hour)

		err = repositories.StoreSession(user.ID, sessionToken, expiryTime)
		if err != nil {
			log.Printf("Failed to store session token: %v", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	} else if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("frontend/templates/sign-in.html")
		if err != nil {
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)

	} else {
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
