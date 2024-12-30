package router

// SetupRoutes configures and returns all application routes
func SetupRoutes() {
	// Register Admin routes
	RegisterAdminRoutes()

	// Register API routes
	RegisterAPIRoutes()

	// Register Newsletter Management routes
	RegisterNewsletterManagementRoutes()

	// Register Blog Management routes
	RegisterBlogManagementRoutes()

	// Register Project Management routes
	RegisterProjectManagementRoutes()

	// Register Analytics Management routes
	RegisterAnalyticsManagementRoutes()
}
