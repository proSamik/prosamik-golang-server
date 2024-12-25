package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
)

// SetupRoutes configures and returns all application routes
func SetupRoutes() {
	// Register routes with middleware
	http.HandleFunc("/readme",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.MarkdownHandler,
			),
		),
	)

	// Register routes for custom repo list
	http.HandleFunc("/repos-list",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleReposList,
			),
		),
	)
}
