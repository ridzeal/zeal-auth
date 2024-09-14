package handler

import (
	"net/http"
	"sso-backend/auth"
	"strings"

	"github.com/labstack/echo/v4"
)

func Protected(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.String(http.StatusUnauthorized, "Missing authorization header")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.String(http.StatusUnauthorized, "Invalid authorization header")
	}

	tokenString := tokenParts[1]
	token, err := auth.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return c.String(http.StatusUnauthorized, "Invalid token")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Access granted to protected resource"})
}
