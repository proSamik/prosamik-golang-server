package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
)

func RegisterAnalyticsManagementRoutes() {
	// Helper function to apply all middlewares
	// Reason: Creates a single point to manage all middleware applications,
	// making it easier to maintain and modify a middleware chain
	withMiddlewares := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(h),
			),
		)
	}

	// Analytics management routes
	// Reason: Using a map makes the route structure clear and reduces repetitive code
	routes := map[string]http.HandlerFunc{
		"/analytics/management": handler.HandleAnalyticsManagement,
		"/analytics/filter":     handler.HandleAnalyticsFilter,
	}

	// Register all routes with middlewares
	// Reason: Single loop to register all routes reduces code duplication
	for path, handlers := range routes {
		http.HandleFunc(path, withMiddlewares(handlers))
	}
}
