package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
)

func RegisterProjectManagementRoutes() {
	// Helper function to apply all middlewares
	withMiddlewares := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(h),
			),
		)
	}

	// Project management routes
	routes := map[string]http.HandlerFunc{
		// Main management route
		"/project/management": handler.HandleProjectManagement,

		// Project search route
		"/project/management/search": handler.HandleProjectSearch,

		// Project add route
		"/project/management/add": handler.HandleProjectAdd,

		// project edit route
		"/project/management/edit/":   handler.HandleProjectEdit,
		"/project/management/update/": handler.HandleProjectUpdate,

		// Project delete rout
		"/project/management/delete/": handler.HandleProjectDelete,

		// Project Cancel-edit
		"/project/management/cancel-edit/": handler.HandleProjectCancelEdit,
	}

	// Register all routes with middlewares
	for path, handlers := range routes {
		http.HandleFunc(path, withMiddlewares(handlers))
	}
}
