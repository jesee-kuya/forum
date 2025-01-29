package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jesee-kuya/forum/backend/database"
	"github.com/jesee-kuya/forum/backend/repositories"
	"github.com/jesee-kuya/forum/backend/routes"
)

func main() {
	db := database.CreateConnection()
	defer db.Close()

	repositories.InitRepo(db)

	router := routes.RegisterRoutes()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server is running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
