package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
)

// RegisterRepoViewRoutes sets up routes for repository view tracking statistics
func RegisterRepoViewRoutes() {
	// Helper function to apply all middlewares
	withMiddlewares := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(h),
			),
		)
	}

	// Repository view routes
	routes := map[string]http.HandlerFunc{
		// Public API route for getting repository views
		"/api/repo/views": handler.GetRepoViewsHandler,

		// Admin routes requiring authentication
		"/api/repo/top":         handler.GetTopRepositoriesHandler,
		"/api/repo/views/date":  handler.GetRepoViewsByDateRangeHandler,
		"/repo/views/analytics": handler.HandleRepoViewAnalytics,
		"/repo/views/filter":    handler.HandleRepoViewFilter,
	}

	// Register all routes with middlewares
	for path, routeHandler := range routes {
		http.HandleFunc(path, withMiddlewares(routeHandler))
	}
}
