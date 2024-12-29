package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
)

func RegisterBlogManagementRoutes() {
	// Helper function to apply all middlewares
	withMiddlewares := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(h),
			),
		)
	}

	// Blog management routes
	routes := map[string]http.HandlerFunc{
		// Main management route
		"/blog/management": handler.HandleBlogManagement,

		// Blog search route
		"/blog/management/search": handler.HandleBlogSearch,

		// Blog add routes
		"/blog/management/add": handler.HandleBlogAdd,

		// Blog edit routes
		"/blog/management/edit/":   handler.HandleBlogEdit,
		"/blog/management/update/": handler.HandleBlogUpdate,

		// Blog delete route
		"/blog/management/delete/": handler.HandleBlogDelete,

		// Cancel-edit
		"/blog/management/cancel-edit/": handler.HandleBlogCancelEdit,
	}

	// Register all routes with middlewares
	for path, handlers := range routes {
		http.HandleFunc(path, withMiddlewares(handlers))
	}
}
