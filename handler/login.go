package handler

import (
	"log"
	"net/http"
	"sso-backend/auth"
	"sso-backend/db"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Query the database for the user's credentials
	var storedCredential string
	var isActive bool
	err := db.DB.QueryRow(`
        SELECT uc.credential, u.is_active
        FROM "z-auth".usercredentials uc
				LEFT JOIN "z-auth".users u ON uc.user_id = u.id
        WHERE uc.auth_method_id = (SELECT id FROM "z-auth".authmethods WHERE method_name = 'email')
        AND uc.user_id = (SELECT id FROM "z-auth".users where username=$1)
    `, email).Scan(&storedCredential, &isActive)
	if err != nil {
		log.Printf("[ERR] %s", err.Error())
		return c.String(http.StatusUnauthorized, "Invalid credentials")
	}

	if !isActive {
		return c.String(http.StatusUnauthorized, "User is not active")
	}

	parts := strings.Split(storedCredential, "::")
	if len(parts) != 2 {
		return c.String(http.StatusInternalServerError, "Invalid stored credential format")
	}

	storedEmail, storedHash := parts[0], parts[1]

	// Compare the stored hash with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil || storedEmail != email {
		return c.String(http.StatusUnauthorized, "Invalid credentials")
	}

	// If we get here, the credentials are valid
	token, err := auth.GenerateToken(email)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not generate token")
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
