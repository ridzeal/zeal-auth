package handler

import (
	"database/sql"
	"log"
	"net/http"
	"sso-backend/db"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationRequest struct {
	OrganizationName string `json:"organization_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
}

func Register(c echo.Context) error {
	var req RegistrationRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	defer tx.Rollback()

	// Insert organization
	var orgID string
	err = tx.QueryRow(`
        INSERT INTO "z-auth".organizations (name) 
        VALUES ($1) 
        RETURNING id
    `, req.OrganizationName).Scan(&orgID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create organization")
	}

	// Insert user
	var userID string
	err = tx.QueryRow(`
			INSERT INTO "z-auth".users (username) 
			VALUES ($1) 
			RETURNING id
	`, req.Email).Scan(&userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create user")
	}

	// Link user to organization
	_, err = tx.Exec(`
			INSERT INTO "z-auth".userorganizations (user_id, organization_id)
			VALUES ($1, $2)
	`, userID, orgID)
	if err != nil {
		log.Printf("tx.Exec Error: %s", err.Error())
		return c.String(http.StatusInternalServerError, "Failed to link user to organization")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to hash password")
	}

	// Get auth method id for password
	authMethodID, err, errMessage := getAuthMethod(tx)
	if err != nil {
		return c.String(http.StatusInternalServerError, errMessage)
	}

	// Insert user credentials
	_, err = tx.Exec(`
        INSERT INTO "z-auth".usercredentials (credential, auth_method_id, user_id) 
        VALUES ($1, $2, $3)
    `, req.Email+"::"+string(hashedPassword), authMethodID, userID)
	if err != nil {
		log.Printf("[ERR] %s", err.Error())
		return c.String(http.StatusInternalServerError, "Failed to store credentials")
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to complete registration")
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Registration successful"})
}

func getAuthMethod(tx *sql.Tx) (authMethodID int, err error, errMessage string) {
	err = tx.QueryRow(`
			SELECT id FROM "z-auth".authmethods WHERE method_name = 'email'
	`).Scan(&authMethodID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert new auth method if it doesn't exist
			err = tx.QueryRow(`
							INSERT INTO "z-auth".authmethods (method_name)
							VALUES ('email')
							RETURNING id
					`).Scan(&authMethodID)
			if err != nil {
				errMessage = "Failed to create auth method"
				return
			}
		} else {
			errMessage = "Failed to get auth method"
			return
		}
	}
	return
}
