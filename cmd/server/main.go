package main

import (
	"log"
	"net/http"
	"prosamik-backend/internal/database"
	"prosamik-backend/internal/router"
)

func main() {

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
