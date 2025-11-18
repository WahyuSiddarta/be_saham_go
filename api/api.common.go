package api

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var Logger *zerolog.Logger

type API struct {
	Router       *echo.Echo
	ServerIP     string
	ServerStatus string
}

// parseUserID parses and validates user ID from path parameter
func parseUserID(c echo.Context, paramName string) (int, error) {
	userIDParam := c.Param(paramName)
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
