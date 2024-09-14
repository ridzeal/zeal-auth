package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Activation(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{})
}
