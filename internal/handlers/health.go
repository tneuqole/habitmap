package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetHealth serves as a basic health check
func GetHealth(c echo.Context) error {
	return c.String(http.StatusOK, "healthy")
}
