package router

import (
	"github.com/WahyuSiddarta/be_saham_go/api"
	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

// setupProtectedRoutes configures routes that require authentication
func (r *Router) setupProtectedRoutes(apiGroup *echo.Group) {
	rprotected := apiGroup.Group("/protected")

	// Add authentication middleware to protected routes
	rprotected.Use(middleware.RequireAuth())

	// Initialize auth handlers
	authHandlers := api.NewAuthHandlers()

	// Setup user routes
	r.setupUserRoutes(rprotected, authHandlers)

	// Setup admin routes
	r.setupAdminRoutes(rprotected, authHandlers)
}

// setupUserRoutes configures user profile and auth endpoints
func (r *Router) setupUserRoutes(rprotected *echo.Group, authHandlers *api.AuthHandlers) {
	userGroup := rprotected.Group("/user")

	// Get current user profile - accessible at /api/protected/user/profile
	userGroup.GET("/profile", authHandlers.GetProfile)

	// User logout - accessible at /api/protected/user/logout
	userGroup.POST("/logout", authHandlers.Logout)
}

// setupAdminRoutes configures admin routes (admin authentication required)
func (r *Router) setupAdminRoutes(rprotected *echo.Group, authHandlers *api.AuthHandlers) {
	adminGroup := rprotected.Group("/admin")
	adminGroup.Use(middleware.AdminRequired())

	// User management endpoints
	usersGroup := adminGroup.Group("/users")

	// Get all users with pagination and filters - accessible at /api/protected/admin/users
	usersGroup.GET("", authHandlers.GetAllUsers, validator.ValidateQuery(&validator.GetUsersQuery{}))

	// Update user level - accessible at /api/protected/admin/users/:id/level
	usersGroup.PUT("/:id/level", authHandlers.UpdateUserLevel, validator.ValidateRequest(&validator.UpdateUserLevelRequest{}))

	// Update user status - accessible at /api/protected/admin/users/:id/status
	usersGroup.PUT("/:id/status", authHandlers.UpdateUserStatus, validator.ValidateRequest(&validator.UpdateUserStatusRequest{}))

	// Get expired users - accessible at /api/protected/admin/users/expired
	usersGroup.GET("/expired", authHandlers.GetExpiredUsers)

	// Downgrade expired users - accessible at /api/protected/admin/users/downgrade-expired
	usersGroup.POST("/downgrade-expired", authHandlers.DowngradeExpiredUsers)
}
