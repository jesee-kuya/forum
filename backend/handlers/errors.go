package handlers

import (
	"net/http"
	"text/template"
	"strconv"
)

type Message struct {
	Code string
	ErrMessage string
}

func ErrorHandler(w http.ResponseWriter, errval string, statusCode int) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, http.StatusInternalServerError)
	}

	code := strconv.Itoa(statusCode)

	var data = message {
		Code: code,
		ErrMessage: errval,
	}

	w.WriteHeader(statusCode)

	err = tmpl.Execute(w, data)
}