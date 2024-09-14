package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sso-backend/auth"
	"sso-backend/db"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

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
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !isActive {
		http.Error(w, "User is not active", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(storedCredential, "::")
	if len(parts) != 2 {
		http.Error(w, "Invalid stored credential format", http.StatusInternalServerError)
		return
	}

	storedEmail, storedHash := parts[0], parts[1]

	// Compare the stored hash with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil || storedEmail != email {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// If we get here, the credentials are valid
	token, err := auth.GenerateToken(email)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
