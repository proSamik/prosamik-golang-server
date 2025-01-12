package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"prosamik-backend/internal/cache"
	"prosamik-backend/internal/database"
	"prosamik-backend/internal/router"
)

func main() {
	// Check and load environment variables in development mode
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Warning: Error loading .env file")
		} else {
			fmt.Println(".env file loaded successfully")
		}
	}

	// Setup routes
	router.SetupRoutes()

	// Initialize database connection and schema
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}

	// Initialize Redis connection
	if err := cache.InitRedis(); err != nil {
		log.Fatal(err)
	}

	// Start server
	port := ":10000"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
