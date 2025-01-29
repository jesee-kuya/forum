package handler

import (
	"fmt"
	"net/http"
	"text/template"
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

	// Method
	if r.Method == http.MethodGet {
		fmt.Println("OK: ", http.StatusOK)
	}

	// template
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

