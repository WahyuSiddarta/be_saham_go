package middleware

import (
	"net/http"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware provides authentication middleware (placeholder for future implementation)
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: Implement JWT or session-based authentication
			// For now, this is a placeholder that allows all requests

			// Example of what authentication might look like:
			// authHeader := c.Request().Header.Get("Authorization")
			// if authHeader == "" {
			//     return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			//         "status":  "ERR",
			//         "code":    "UNAUTHORIZED",
			//         "message": "Authorization header required",
			//     })
			// }

			// Validate token here...
			// if !isValidToken(authHeader) {
			//     return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			//         "status":  "ERR",
			//         "code":    "INVALID_TOKEN",
			//         "message": "Invalid or expired token",
			//     })
			// }

			return next(c)
		}
	}
}

// RequireAuth returns a middleware that requires authentication
func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: Implement actual authentication check
			// This is a placeholder that returns unauthorized for now

			if Logger != nil {
				Logger.Warn().Str("path", c.Request().URL.Path).Msg("Authentication required but not implemented")
			}

			return helper.ErrorResponse(c, http.StatusUnauthorized, "Sistem autentikasi belum diimplementasi", "Endpoint ini memerlukan autentikasi ketika sistem auth telah diaktifkan")
		}
	}
}
