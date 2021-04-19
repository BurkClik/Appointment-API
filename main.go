package main

import (
	"github.com/BurkClik/Appointment-API/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	// Echo instance
	e := echo.New()

	routes.SetupRoutes(e)
}
