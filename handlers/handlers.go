package handlers

import (
	"net/http"
	"sso-backend/auth"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// In a real application, you would validate the credentials against a database
	if username == "admin" && password == "password" {
		token, err := auth.GenerateToken(username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
		}
		return c.JSON(http.StatusOK, map[string]string{"token": token})
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
}

func Protected(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
	}

	tokenString := authHeader[7:] // Remove "Bearer " prefix
	token, err := auth.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Access granted to protected resource"})
}