package controller

import (
	"errors"

	"github.com/fleimkeipa/tickets-api/pkg"

	"github.com/labstack/echo/v4"
)

// HandleEchoError handles errors that occur within the Echo framework.
func HandleEchoError(c echo.Context, err error) error {
	var pe *pkg.Error

	if errors.As(err, &pe) {
		return c.JSON(pe.StatusCode(), pe.Message())
	} else {
		// Log the error
		return c.JSON(500, "Internal Server Error")
	}
}
