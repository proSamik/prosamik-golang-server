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

	// Register route for feedback form
	http.HandleFunc("/feedback",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleFeedback,
			),
		),
	)

	// Register route for newsletter subscription
	http.HandleFunc("/newsletter",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleNewsletter,
			),
		),
	)

	// Root route handler
	http.HandleFunc("/", handler.HandleRoot)

	// Admin routes
	http.HandleFunc("/samik/login",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleAdminLogin,
			),
		),
	)

	http.HandleFunc("/samik/logout",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleLogout,
				),
			),
		),
	)

	http.HandleFunc("/samik",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleAdminDashboard,
				),
			),
		),
	)
}
