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
	util.Init()
	defer util.DB.Close()

	port, err := util.ValidatePort()
	if err != nil {
		log.Fatalf("Error validating port: %v", err)
		return
	}
	r := route.InitRoutes()

	server := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	url := fmt.Sprintf("https://localhost:%v", port)

	go http.ListenAndServe("+" + port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url+r.RequestURI, http.StatusMovedPermanently)
	}))

	log.Printf("Server started at https://localhost%s\n", port)
	if err = server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
