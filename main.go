package main

import (
	"log"
	"net/http"
	"sso-backend/api"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()
	api.Setup(e)
	s := http.Server{
		Addr:    ":3000",
		Handler: e,
	}
	// Start the server
	log.Println("Server starting on http://localhost:3000")
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
