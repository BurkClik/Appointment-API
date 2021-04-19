package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BurkClik/Appointment-API/helper"
	"github.com/BurkClik/Appointment-API/model"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

var database = helper.ConnectDB()

// Hello is
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

// Signup : This is a sign up for non-doctors people
func Signup(c echo.Context) error {

	var user model.User

	_ = json.NewDecoder(c.Request().Body).Decode(&user)

	// Validate
	if user.Mail == "" || user.Name == "" || user.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or password"}
	} else if user.City == "" || user.Gender == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid city or gender"}
	}

	// Save user
	insertResult, err := database.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return c.JSON(http.StatusCreated, user)
}


// SignupDoctor : This is a sign up for doctors
func SignupDoctor(c echo.Context) error {

	var doctor model.Doctor

	_ = json.NewDecoder(c.Request().Body).Decode(&doctor)

	// Validate
	if doctor.Mail == "" || doctor.Name == "" || doctor.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or password"}
	}

	// Save doctor
	insertResult, err := database.Collection("doctors").InsertOne(context.TODO(), doctor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return c.JSON(http.StatusCreated, doctor)
}