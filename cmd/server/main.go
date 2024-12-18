package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"prosamik-backend/internal/auth"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
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
	http.HandleFunc("/readme",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleReadmeRequest,
			),
		),
	)

	// Start server
	port := ":10000"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
