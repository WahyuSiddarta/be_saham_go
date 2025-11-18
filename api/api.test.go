package api

import "github.com/labstack/echo/v4"

func (api *API) Test(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"status":  "success",
		"message": "API is working",
	})
}
