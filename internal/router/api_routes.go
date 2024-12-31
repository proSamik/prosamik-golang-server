package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
	"time"
)

func RegisterAPIRoutes() {
	// Rate limiter for feedback and newsletter
	// Reason: Initialize rate limiter once to be used across multiple routes
	rateLimiter := middleware.NewRateLimiter(60, time.Minute)

	// Helper function for standard middleware chain
	// Reason: Creates a reusable middleware stack for regular routes
	withStandardMiddlewares := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.CORSMiddleware(
			middleware.LoggingMiddleware(h),
		)
	}

	// Helper function for rate-limited middleware chain
	// Reason: Creates a separate middleware stack for routes that need rate limiting
	withRateLimitedMiddlewares := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				rateLimiter.RateLimitMiddleware(h),
			),
		)
	}

	// Standard routes without rate limiting
	// Reason: Group similar routes together for better organization
	standardRoutes := map[string]http.HandlerFunc{
		"/blogs":     handler.HandleBlogsList,
		"/projects":  handler.HandleProjectsList,
		"/md":        handler.MarkdownHandler,
		"/analytics": handler.HandleAnalytics,
	}

	// Rate-limited routes
	// Reason: Separate routes that need rate limiting for clarity
	rateLimitedRoutes := map[string]http.HandlerFunc{
		"/feedback":   handler.HandleFeedback,
		"/newsletter": handler.HandleNewsletterSignup,
	}

	// Register standard routes
	// Reason: Apply standard middleware stack to regular routes
	for path, apiHandlers := range standardRoutes {
		http.HandleFunc(path, withStandardMiddlewares(apiHandlers))
	}

	// Register rate-limited routes
	// Reason: Apply rate-limited middleware stack to routes that need it
	for path, handlers := range rateLimitedRoutes {
		http.HandleFunc(path, withRateLimitedMiddlewares(handlers))
	}
}
