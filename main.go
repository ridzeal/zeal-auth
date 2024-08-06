package main

import (
	"fmt"
	"net/http"
	"sso-backend/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.POST("/login", handlers.Login)
	e.GET("/protected", handlers.Protected)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}
