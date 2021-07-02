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
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

var database = helper.ConnectDB()

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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

	user.Password, _ = HashPassword(user.Password)

	// Save user
	insertResult, err := database.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return c.JSON(http.StatusOK, "Bu email ile kayıtlı hesap var")
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	//----------
	// JWT
	//----------

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	user.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}

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

// HospitalList :
func HospitalList(c echo.Context) error {
	cursor, err := database.Collection("hospitals").Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var hospitals []bson.M
	if err = cursor.All(context.TODO(), &hospitals); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, hospitals)
}

// TopRatedDoctors :
func TopRatedDoctors(c echo.Context) error {

	doctorCount := c.QueryParam("count")
	i, err := strconv.ParseInt(doctorCount, 10, 32)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"doctor.vote_rate", -1}})

	limit := options.Find()
	limit.SetLimit(i)

	cursor, err := database.Collection("users").Find(context.TODO(), bson.M{"is_doctor": true}, findOptions, limit)
	if err != nil {
		log.Fatal(err)
	}

	var doctors []bson.M
	if err = cursor.All(context.TODO(), &doctors); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, doctors)
}

// SearchDoctor :
func SearchDoctor(c echo.Context) error {

	doctor := c.QueryParam("doctor")

	var requestError model.Error

	isDoctor := bson.M{"is_doctor": true}
	queryResult := bson.M{
		"$or": []bson.M{bson.M{"name": primitive.Regex{Pattern: doctor, Options: ""}},
			bson.M{"doctor.department": primitive.Regex{Pattern: doctor, Options: ""}}}}

	cursor, err := database.Collection("users").Find(context.TODO(),
		bson.M{"$and": []bson.M{isDoctor, queryResult}})
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	if results == nil {
		requestError.Code = http.StatusOK
		requestError.Message = "Aradığınız veri bulunamadı"
		return c.JSON(http.StatusOK, requestError)
	}

	return c.JSON(http.StatusOK, results)
}

// SearchHospital :
func SearchHospital(c echo.Context) error {

	hospital := c.QueryParam("hospital")

	var requestError model.Error

	cursor, err := database.Collection("users").Find(context.TODO(),
		bson.M{"doctor.hospital": primitive.Regex{Pattern: hospital, Options: ""}})
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	if results == nil {
		requestError.Code = http.StatusOK
		requestError.Message = "Aradığınız veri bulunamadı"
		return c.JSON(http.StatusOK, requestError)
	}

	return c.JSON(http.StatusOK, results)
}

// UserDetail :
func UserDetail(c echo.Context) error {
	mail := c.QueryParam("mail")

	filterCursor, err := database.Collection("users").Find(context.TODO(), bson.M{"mail": mail})

	if err != nil {
		log.Fatal(err)
	}

	var userFiltered []bson.M
	if err = filterCursor.All(context.TODO(), &userFiltered); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, userFiltered)
}

func DoctorDetail(c echo.Context) error {
	id := c.QueryParam("_id")

	var requestError model.Error

	log.Println("id -> " + id)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		requestError.Code = http.StatusOK
		requestError.Message = err.Error()
		return c.JSON(http.StatusOK, requestError)
	}

	var doctor bson.M

	if err = database.Collection("users").FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&doctor); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, doctor)
}

func GetAppointment(c echo.Context) error {
	userId := c.QueryParam("_id")

	var appointment model.Appointment

	_ = json.NewDecoder(c.Request().Body).Decode(&appointment)

	log.Println("Saat " + appointment.Hour)

	id, _ := primitive.ObjectIDFromHex(userId)

	who := bson.M{"_id": id}

	updateAppointment := bson.D{{"$push", bson.D{{"appointment", appointment}}}}

	result, err := database.Collection("users").UpdateOne(context.TODO(), who, updateAppointment)

	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, result)
}

func MakeReview(c echo.Context) error {
	userId := c.QueryParam("_id")

	var review model.Review

	_ = json.NewDecoder(c.Request().Body).Decode(&review)

	id, _ := primitive.ObjectIDFromHex(userId)

	who := bson.M{"_id": id}

	updateReview := bson.D{{"$push", bson.D{{"review", review}}}}

	result, err := database.Collection("users").UpdateOne(context.TODO(), who, updateReview)

	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, result)
}
