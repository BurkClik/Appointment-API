package routes

import (
	"github.com/BurkClik/Appointment-API/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(e *echo.Echo) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handler.Hello)
	e.GET("/doctors", handler.DoctorList)
	e.GET("/search/:query", handler.Search)

	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
