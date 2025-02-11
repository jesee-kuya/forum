package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"text/template"

	"github.com/jesee-kuya/forum/backend/models"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/util"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if r.URL.Path != "/sign-up" {
		util.ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed parsing form: %v\n", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}

		user.Username = strings.TrimSpace(r.PostFormValue("username"))
		user.Email = strings.TrimSpace(r.PostFormValue("email"))
		user.Password = strings.TrimSpace(r.PostFormValue("password"))
		user.ConfirmedPassword = strings.TrimSpace(r.PostFormValue("confirmed-password"))

		err = util.ValidateFormFields(user.Username, user.Email, user.Password)
		if err != nil {
			log.Printf("Invalid form values from user: %v\n", err)
			response := Response{Success: false}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if !reflect.DeepEqual(user.Password, user.ConfirmedPassword) {
			log.Printf("User password %q and confirmed password %q do not match.\n", user.Password, user.ConfirmedPassword)
			response := Response{Success: false}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		hashed, err := util.PasswordEncrypt([]byte(user.Password), 10)
		if err != nil {
			log.Printf("Failed encrypting password: %v\n", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}

		_, err = repositories.InsertRecord(util.DB, "tblUsers", []string{"username", "email", "user_password"}, user.Username, user.Email, string(hashed))
		if err != nil {
			log.Println("Error adding user:", err)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	} else if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("frontend/templates/sign-up.html")
		if err != nil {
			log.Println("failed parsing files:", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else {
		log.Println("Method not allowed", r.Method)
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
