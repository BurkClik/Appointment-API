package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BurkClik/Appointment-API/helper"
	"github.com/BurkClik/Appointment-API/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
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
	}

	// Save user
	insertResult, err := database.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return c.JSON(http.StatusOK, "Bu email ile kayıtlı hesap var")
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return c.JSON(http.StatusCreated, user)
}

// Login :
func Login(c echo.Context) (err error) {
	// Bind
	u := new(model.User)
	if err = c.Bind(u); err != nil {
		return
	}

	_ = json.NewDecoder(c.Request().Body).Decode(&u)

	filter := bson.M{"mail": u.Mail, "password": u.Password}

	// Find user
	if err = database.Collection("users").FindOne(context.TODO(), filter).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid email or password"}
		}
	}

	//---------
	// JWT
	//---------

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	u.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}

	u.Password = "" // Don't send password
	return c.JSON(http.StatusOK, u)
}

// DoctorList :
func DoctorList(c echo.Context) error {

	cursor, err := database.Collection("users").Find(context.TODO(), bson.M{"is_doctor": true})
	if err != nil {
		log.Fatal(err)
	}

	var doctors []bson.M
	if err = cursor.All(context.TODO(), &doctors); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, doctors)
}

// Search :
func Search(c echo.Context) error {

	query := c.Param("query")

	isDoctor := bson.M{"is_doctor": true}
	queryResult := bson.M{
		"$or": []bson.M{bson.M{"name": primitive.Regex{Pattern: query, Options: ""}},
			bson.M{"doctor.department": primitive.Regex{Pattern: query, Options: ""}}}}

	cursor, err := database.Collection("users").Find(context.TODO(),
		bson.M{"$and": []bson.M{isDoctor, queryResult}})
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, results)
}
