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
	e.GET("/toprated/doctors", handler.TopRatedDoctors)
	e.GET("/search/doctor", handler.SearchDoctor)
	e.GET("/search/hospital", handler.SearchHospital)
	e.GET("/user", handler.UserDetail)
	e.GET("/doctor", handler.DoctorDetail)
	e.GET("/hospital", handler.HospitalList)

	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)
	e.POST("/appointment", handler.GetAppointment)
	e.POST("/review", handler.MakeReview)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
