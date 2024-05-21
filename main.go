package main

import (
	"go-rate-limiter/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
