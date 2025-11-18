package router

import (
	"github.com/WahyuSiddarta/be_saham_go/api"
	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

// setupProtectedRoutes configures routes that require authentication
func (r *Router) setupProtectedRoutes(apiGroup *echo.Group) {
	// Initialize auth handlers
	authHandlers := api.NewAuthHandlers()

	userGroup := apiGroup.Group("/users")
	// Add authentication middleware to protected routes
	userGroup.Use(middleware.RequireAuth())

	// Get current user profile - accessible at /api/protected/user/profile
	userGroup.GET("/profile", authHandlers.GetProfile)

	// Setup portfolio routes (includes PnL)
	setupPortfolioRoutes(userGroup)

	// Setup admin routes
	r.setupAdminRoutes(apiGroup, authHandlers)

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

}

// setupPortfolioRoutes configures portfolio cash routes
func setupPortfolioRoutes(rprotected *echo.Group) {
	portfolioGroup := rprotected.Group("/portfolio")

	// Initialize portfolio handlers
	portfolioRepo := models.NewPortfolioCashRepository()
	portfolioHandlers := api.NewPortfolioCashHandlers(portfolioRepo)

	// Create cash portfolio - accessible at /api/protected/portfolio/cash
	portfolioGroup.POST("/cash", portfolioHandlers.CreateCashPortfolio, validator.ValidateRequest(&validator.CreatePortfolioCashRequest{}))

	// Get all cash portfolios - accessible at /api/protected/portfolio/cash
	portfolioGroup.GET("/cash", portfolioHandlers.GetMyCashPortfolios)

	// Update cash portfolio - accessible at /api/protected/portfolio/cash/:id
	portfolioGroup.PUT("/cash/:id", portfolioHandlers.UpdateCashPortfolio, validator.ValidateRequest(&validator.UpdatePortfolioCashRequest{}))

	// Delete cash portfolio - accessible at /api/protected/portfolio/cash/:id
	portfolioGroup.DELETE("/cash/:id", portfolioHandlers.DeleteCashPortfolio)

	// Move asset between portfolios - accessible at /api/protected/portfolio/move-asset
	portfolioGroup.POST("/move-asset", portfolioHandlers.MoveAsset, validator.ValidateRequest(&validator.MoveAssetRequest{}))

	// Realize cash portfolio - accessible at /api/protected/portfolio/realize
	portfolioGroup.POST("/realize", portfolioHandlers.RealizeCashPortfolio, validator.ValidateRequest(&validator.RealizeCashPortfolioRequest{}))

	// PnL sub-routes under portfolio
	pnlGroup := portfolioGroup.Group("/cash/pnl")
	pnlGroup.GET("", portfolioHandlers.GetPnlRealizedCash)
	pnlGroup.GET("/portfolio/:portfolioId", portfolioHandlers.GetPnlByPortfolioCashID)
	pnlGroup.GET("/:id", portfolioHandlers.GetPnlById)
	pnlGroup.POST("", portfolioHandlers.CreatePnlRealizedCash, validator.ValidateRequest(&validator.CreatePnlRealizedCashRequest{}))
	pnlGroup.PUT("/:id", portfolioHandlers.UpdatePnlRealizedCash, validator.ValidateRequest(&validator.UpdatePnlRealizedCashRequest{}))
	pnlGroup.DELETE("/:id", portfolioHandlers.DeletePnlRealizedCash)
}
