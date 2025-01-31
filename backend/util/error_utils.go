package util

import (
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type Message struct {
	Code       string
	ErrMessage string
}

func ErrorHandler(w http.ResponseWriter, errval string, statusCode int) {
	tmpl, err := template.ParseFiles("frontend/templates/error.html")
	if err != nil {
		log.Printf("Failed to load error template: %v", err)
		http.Error(w, errval, statusCode)
		return
	}

	code := strconv.Itoa(statusCode)

	data := Message{
		Code:       code,
		ErrMessage: errval,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing the template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
