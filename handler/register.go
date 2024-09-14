package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sso-backend/db"

	"golang.org/x/crypto/bcrypt"
)

type RegistrationRequest struct {
	OrganizationName string `json:"organization_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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
		http.Error(w, "Failed to create organization", http.StatusInternalServerError)
		return
	}

	// Insert user
	var userID string
	err = tx.QueryRow(`
			INSERT INTO "z-auth".users (username) 
			VALUES ($1) 
			RETURNING id
	`, req.Email).Scan(&userID)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Link user to organization
	_, err = tx.Exec(`
			INSERT INTO "z-auth".userorganizations (user_id, organization_id)
			VALUES ($1, $2)
	`, userID, orgID)
	if err != nil {
		log.Printf("tx.Exec Error: %s", err.Error())
		http.Error(w, "Failed to link user to organization", http.StatusInternalServerError)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Get auth method id for password
	authMethodID, err, errMessage := getAuthMethod(tx)
	if err != nil {
		http.Error(w, errMessage, http.StatusInternalServerError)
	}

	// Insert user credentials
	_, err = tx.Exec(`
        INSERT INTO "z-auth".usercredentials (credential, auth_method_id, user_id) 
        VALUES ($1, $2, $3)
    `, req.Email+"::"+string(hashedPassword), authMethodID, userID)
	if err != nil {
		log.Printf("[ERR] %s", err.Error())
		http.Error(w, "Failed to store credentials", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Failed to complete registration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
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
