package main

import (
	"log"
	"net/http"
	"prosamik-backend/internal/router"
)

func main() {

	// Setup routes
	router.SetupRoutes()

	// Start server
	port := ":10000"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
