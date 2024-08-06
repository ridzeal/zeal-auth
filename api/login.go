package api

import (
    "encoding/json"
    "net/http"
    "sso-backend/auth"
)

func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    if username == "admin" && password == "password" {
        token, err := auth.GenerateToken(username)
        if err != nil {
            http.Error(w, "Could not generate token", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"token": token})
    } else {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
    }
}