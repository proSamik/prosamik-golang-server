package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"prosamik-backend/internal/auth"
	"prosamik-backend/internal/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Validate GitHub token
	if err := auth.ValidateGitHubToken(); err != nil {
		log.Fatalf("Authentication error: %v", err)
	}

	// Setup routes
	router.SetupRoutes()

	// Start server
	port := ":10000"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
