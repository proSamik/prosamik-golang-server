package router

// SetupRoutes configures and returns all application routes
func SetupRoutes() {
	// Register API routes
	RegisterAPIRoutes()

	// Register Admin routes
	RegisterAdminRoutes()

	// Register Newsletter Management routes
	RegisterNewsletterManagementRoutes()

	// Register Blog Management routes
	RegisterBlogManagementRoutes()

	// Register Project Management routes
	RegisterProjectManagementRoutes()
}
