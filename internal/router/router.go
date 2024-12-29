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

	// Admin routes
	http.HandleFunc("/login",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleAdminLoginUsingJWT,
			),
		),
	)

	http.HandleFunc("/logout",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleAdminLogout,
				),
			),
		),
	)

	http.HandleFunc("/",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleDashboard,
				),
			),
		),
	)

	// Register routes for blogs
	http.HandleFunc("/blogs",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleBlogsList,
			),
		),
	)

	// Register routes for markdown
	http.HandleFunc("/md",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.MarkdownHandler,
			),
		),
	)

	// New route for projects
	http.HandleFunc("/projects",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleProjectsList,
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

	// Newsletter management route
	http.HandleFunc("/newsletter/management",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleNewsletterManagement,
				),
			),
		),
	)

	// Newsletter delete route
	http.HandleFunc("/newsletter/delete/",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleNewsletterDelete,
				),
			),
		),
	)

}
