package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jesee-kuya/forum/backend/route"
	"github.com/jesee-kuya/forum/backend/util"
	"github.com/joho/godotenv"
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

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")

	if certPath == "" || keyPath == "" {
		log.Fatal("Certificate path or key path not specified")
	}

	server := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
		},
	}

	log.Printf("Server started at https://localhost%s\n", port)
	if err = server.ListenAndServeTLS(certPath, keyPath); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
