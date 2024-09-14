package validation

import (
	"log"
	"net/http"
	"sso-backend/db"
	"sso-backend/state"
)

func ValidateRequest(w http.ResponseWriter, r *http.Request) (err error) {
	if r.Method != http.MethodPost {
		err = &state.Error{
			Message: "Method not allowed",
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return
}
