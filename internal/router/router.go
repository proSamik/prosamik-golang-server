package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
	"time"
)

// SetupRoutes configures and returns all application routes
func SetupRoutes() {

	// Create a rate limiter that allows 5 requests per minute
	rateLimiter := middleware.NewRateLimiter(60, time.Minute)

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
				rateLimiter.RateLimitMiddleware(
					handler.HandleFeedback,
				),
			),
		),
	)

	// Register route for newsletter subscription with rate limiting
	http.HandleFunc("/newsletter",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				rateLimiter.RateLimitMiddleware(
					handler.HandleNewsletter,
				),
			),
		),
	)

	// Root route handler
	http.HandleFunc("/", handler.HandleRoot)

	// Admin routes
	http.HandleFunc("/samik/login",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleAdminLoginUsingJWT,
			),
		),
	)

	http.HandleFunc("/samik/logout",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleAdminLogout,
				),
			),
		),
	)

	http.HandleFunc("/samik",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleDashboard,
				),
			),
		),
	)
}
