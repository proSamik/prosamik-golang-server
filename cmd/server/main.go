package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"prosamik-backend/internal/database"
	"prosamik-backend/internal/router"
)

func main() {

	// Check and load environment variables in development mode
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: Error loading .env file")
		} else {
			log.Println(".env file loaded successfully")
		}
	}

	// Setup routes
	router.SetupRoutes()

	// Initialize database connection and schema
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}

	// Start server
	port := ":10000"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
