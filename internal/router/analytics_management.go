package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
)

func RegisterAnalyticsManagementRoutes() {
	http.HandleFunc("/analytics/management",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleAnalyticsManagement,
				),
			),
		),
	)

	http.HandleFunc("/analytics/filter",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleAnalyticsFilter,
				),
			),
		),
	)
}
