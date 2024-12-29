package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
	"time"
)

func RegisterAPIRoutes() {
	// Rate limiter for feedback and newsletter
	rateLimiter := middleware.NewRateLimiter(60, time.Minute)

	// Blogs route
	http.HandleFunc("/blogs",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleBlogsList,
			),
		),
	)

	// Projects route
	http.HandleFunc("/projects",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleProjectsList,
			),
		),
	)

	// Markdown route
	http.HandleFunc("/md",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.MarkdownHandler,
			),
		),
	)

	// Feedback route with rate limiting
	http.HandleFunc("/feedback",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				rateLimiter.RateLimitMiddleware(
					handler.HandleFeedback,
				),
			),
		),
	)

	// Newsletter subscription with rate limiting
	http.HandleFunc("/newsletter",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				rateLimiter.RateLimitMiddleware(
					handler.HandleNewsletterSignup,
				),
			),
		),
	)
}
