package handler

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/jesee-kuya/forum/backend/models"
)

type Response struct {
	pageTitle string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// URL path
	if r.URL.Path != "/" {
		ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	// Method used
	if r.Method == http.MethodGet {
		fmt.Println("OK: ", http.StatusOK)
	} else {
		ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// template rendering
	tmpl, err := template.ParseFiles("frontend/templates/index.html")
	if err != nil {
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := Response{
		pageTitle: "Home",
	}
	tmpl.Execute(w, data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sign-in" {
		ErrorHandler(w, "Page not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		fmt.Println("OK: ", http.StatusOK)
	} else {
		ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if user.Username == "" || user.Email == "" || user.Password == "" {
		ErrorHandler(w, "Username/Email and Password are required", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("frontend/templates/sign-in.html")
	if err != nil {
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
	}

	data := Response{
		pageTitle: "sign-in",
	}

	tmpl.Execute(w, data)
}
