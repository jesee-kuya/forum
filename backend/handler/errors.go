package handler

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
		http.Error(w, errval, statusCode)
	}

	code := strconv.Itoa(statusCode)

	data := Message{
		Code:       code,
		ErrMessage: errval,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error with the template file", err)
	}
}
