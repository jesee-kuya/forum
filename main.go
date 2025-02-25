package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jesee-kuya/forum/backend/route"
	"github.com/jesee-kuya/forum/backend/util"
)

func main() {
	err := util.LoadEnv(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	util.Init()
	defer util.DB.Close()

	port, err := util.ValidatePort()
	if err != nil {
		log.Fatalf("Error validating port: %v", err)
		return
	}
	router := route.InitRoutes()

	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server started at http://localhost%s\n", port)
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
