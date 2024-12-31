package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
)

func RegisterNewsletterManagementRoutes() {
	// Helper function to apply all middlewares
	// Reason: Creates a single point to manage all middleware applications,
	// making it easier to add/remove middleware for all routes at once
	withMiddlewares := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(h),
			),
		)
	}

	// Newsletter management routes
	// Reason: Using a map makes it easy to see all routes at a glance
	// and reduces the repetitive middleware wrapping code
	routes := map[string]http.HandlerFunc{
		// Main management route
		"/newsletter/management": handler.HandleNewsletterManagement,

		// Newsletter search route
		"/newsletter/search": handler.HandleNewsletterSearch,

		// Newsletter add route
		"/newsletter/add": handler.HandleNewsletterAdd,

		// Newsletter edit routes
		"/newsletter/edit/":   handler.HandleNewsletterEdit,
		"/newsletter/update/": handler.HandleNewsletterUpdate,

		// Newsletter delete route
		"/newsletter/delete/": handler.HandleNewsletterDelete,

		// Cancel-edit
		"/newsletter/cancel-edit/": handler.HandleNewsletterCancelEdit,
	}

	// Register all routes with middlewares
	// Reason: Single loop to register all routes reduces code duplication
	// and makes it less prone to errors
	for path, handlers := range routes {
		http.HandleFunc(path, withMiddlewares(handlers))
	}
}
