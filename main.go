package main

import (
	"log"
	"net/http"
	"sso-backend/db"
	"sso-backend/handler"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Setup(e *echo.Echo) {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", handler.Login)
	e.GET("/protected", handler.Protected)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	e := echo.New()
	Setup(e)
	e.ServeHTTP(w, r)
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()
	Setup(e)
	s := http.Server{
		Addr:    ":3000",
		Handler: e,
		//ReadTimeout: 30 * time.Second, // customize http.Server timeouts
	}
	// Start the server
	log.Println("Server starting on http://localhost:3000")
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
