package main

import (
	"log"
	"net/http"

	"github.com/jesee-kuya/forum/backend/handler"
	"github.com/jesee-kuya/forum/backend/util"
)

func main() {
	util.Init()

	fs := http.FileServer(http.Dir("./frontend/static"))
	http.Handle("/frontend/static/", http.StripPrefix("/frontend/static/", fs))

	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/sign-in", handler.LoginHandler)
	http.HandleFunc("/sign-up", handler.SignupHandler)
	http.HandleFunc("/upload", handler.UploadMedia)

	port, err := util.ValidatePort()
	if err != nil {
		log.Printf("Error validating port: %v\n", err)
		return
	}
	log.Printf("Server started at http://localhost%s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
