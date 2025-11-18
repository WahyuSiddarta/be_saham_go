package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var Logger *zerolog.Logger

type API struct {
	Router       *echo.Echo
	ServerIP     string
	ServerStatus string
}
