package main

import (
	"github.com/BurkClik/Appointment-API/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())


	// Routes
	e.GET("/", handler.Hello)

	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)
	e.POST("/signup/doctor", handler.SignupDoctor)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
