package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
)

func RegisterNewsletterManagementRoutes() {
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

	// Newsletter search route
	http.HandleFunc("/newsletter/search",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleNewsletterSearch,
				),
			),
		),
	)

	// Newsletter add route
	http.HandleFunc("/newsletter/add",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleNewsletterAdd,
				),
			),
		),
	)

	// Newsletter update route
	http.HandleFunc("/newsletter/update/",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleNewsletterUpdate,
				),
			),
		),
	)

	// Newsletter edit route
	http.HandleFunc("/newsletter/edit/",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleNewsletterEdit,
				),
			),
		),
	)

	// Newsletter cancel edit route
	http.HandleFunc("/newsletter/cancel-edit/",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleNewsletterCancelEdit,
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
