package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/jesee-kuya/forum/backend/route"
	"github.com/jesee-kuya/forum/backend/util"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	util.Init()
	defer util.DB.Close()

	port, err := util.ValidatePort()
	if err != nil {
		log.Fatalf("Error validating port: %v", err)
		return
	}
	router := route.InitRoutes()

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("localhost"),
		Cache:      autocert.DirCache("certs"),
	}

	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))

	log.Printf("Server started at http://localhost%s\n", port)
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
