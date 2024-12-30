package router

import (
	"net/http"
	"prosamik-backend/internal/handler"
	"prosamik-backend/internal/middleware"
)

func RegisterAdminRoutes() {
	// Add this FIRST - before other routes
	// Static file server
	fs := http.FileServer(http.Dir("static")) // Assuming your static folder is at project root
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Login route
	http.HandleFunc("/login",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				handler.HandleAdminLoginUsingJWT,
			),
		),
	)

	// Logout route
	http.HandleFunc("/logout",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					handler.HandleAdminLogout,
				),
			),
		),
	)

	// Combined dashboard and catch-all route
	http.HandleFunc("/",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthMiddleware(
					func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Path != "/" { // Catch unmatched routes
							http.Redirect(w, r, "/", http.StatusFound)
							return
						}
						// Dashboard handler logic
						handler.HandleDashboard(w, r)
					},
				),
			),
		),
	)
}
