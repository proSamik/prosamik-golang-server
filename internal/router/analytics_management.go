package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
)

func RegisterAnalyticsManagementRoutes() {
	withMiddlewares := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(h),
			),
		)
	}

	routes := map[string]http.HandlerFunc{
		"/analytics/management": handler.HandleAnalyticsManagement,
		"/analytics/filter":     handler.HandleAnalyticsFilter,
		"/analytics/cache":      handler.HandleCacheMonitoring,
	}

	for path, handlers := range routes {
		http.HandleFunc(path, withMiddlewares(handlers))
	}
}
