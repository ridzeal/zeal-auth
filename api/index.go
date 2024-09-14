package api

import (
	"log"
	"net/http"
	"sso-backend/db"
	"sso-backend/handler"

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
	e.POST("/register", handler.Register)
	e.POST("/activation", handler.Activation)
	e.GET("/protected", handler.Protected)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	e := echo.New()
	Setup(e)
	e.ServeHTTP(w, r)
}
