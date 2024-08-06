package api

import (
    "encoding/json"
    "net/http"
    "strings"
    "sso-backend/auth"
)

func Protected(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Missing authorization header", http.StatusUnauthorized)
        return
    }

    tokenParts := strings.Split(authHeader, " ")
    if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
        http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
        return
    }

    tokenString := tokenParts[1]
    token, err := auth.ValidateToken(tokenString)
    if err != nil || !token.Valid {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Access granted to protected resource"})
}