package main

import (
	"log"
	"net/http"
	"sso-backend/api"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Set up the routes
	http.HandleFunc("/login", api.Login)
	http.HandleFunc("/register", api.Register)
	http.HandleFunc("/activation", api.Activation)
	http.HandleFunc("/protected", api.Protected)

	// Start the server
	log.Println("Server starting on http://localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
