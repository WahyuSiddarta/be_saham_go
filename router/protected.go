package router

import (
	"net/http"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/labstack/echo/v4"
)

// setupProtectedRoutes configures routes that require authentication
func (r *Router) setupProtectedRoutes(apiGroup *echo.Group) {
	rprotected := apiGroup.Group("/protected")

	// Add authentication middleware to protected routes
	rprotected.Use(middleware.RequireAuth())

	// Example protected endpoint - accessible at /api/protected/profile
	rprotected.GET("/profile", func(c echo.Context) error {
		profileData := map[string]interface{}{
			"user_id": 1,
			"email":   "user@example.com",
			"role":    "user",
		}
		return helper.JsonResponse(c, http.StatusOK, profileData)
	})

	// Example: Stock endpoints - these will be accessible at /api/protected/stocks/*
	stockGroup := rprotected.Group("/stocks")
	stockGroup.GET("", func(c echo.Context) error {
		stocksData := []interface{}{
			map[string]interface{}{
				"id":     1,
				"symbol": "BBCA",
				"name":   "Bank Central Asia",
				"price":  10000,
			},
			map[string]interface{}{
				"id":     2,
				"symbol": "BBRI",
				"name":   "Bank Rakyat Indonesia",
				"price":  5000,
			},
		}
		return helper.JsonResponse(c, http.StatusOK, stocksData)
	})
	stockGroup.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		stockData := map[string]interface{}{
			"id":     id,
			"symbol": "BBCA",
			"name":   "Bank Central Asia",
			"price":  10000,
			"change": "+2.5%",
		}
		return helper.JsonResponse(c, http.StatusOK, stockData)
	})
}
