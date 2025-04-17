package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "sample-service/docs"
)

// RegisterSwaggerRoutes registers the swagger routes
func RegisterSwaggerRoutes(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
