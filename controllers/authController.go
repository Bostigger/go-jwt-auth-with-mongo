package controllers

import (
	"context"
	"fmt"
	"github.com/bostigger/go-jwt-mongo/helpers"
	"github.com/bostigger/go-jwt-mongo/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func HashPassword(userPassword string) string {
	password, err := bcrypt.GenerateFromPassword([]byte(userPassword), 14)
	if err != nil {
		return ""
	}
	return string(password)
}

func VerifyPassword(userPassword string, enteredPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(enteredPassword), []byte(userPassword))
	match := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("Password is not correct")
		match = false
	}
	return match, msg

}

func CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var validate = validator.New()
	validateErr := validate.Struct(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validateErr})
		return
	}

	checkEmail, err := usersCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		return
	}
	checkPhone, err := usersCollection.CountDocuments(ctx, bson.M{"phoneNumber": user.PhoneNumber})
	if err != nil {
		return
	}

	if checkEmail > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Email already exists"})
		return
	}

	if checkPhone > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Phone Number already exists"})
		return
	}
	password := HashPassword(user.Password)
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.UserId = user.ID.Hex()
	token, refreshToken, _ := helpers.GenerateAllToken(user.Email, user.Username, user.UserId, user.UserType)
	user.Token = token
	user.RefreshToken = refreshToken
	user.Password = password

	res, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(res)

	c.JSON(http.StatusOK, user)
}

func LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	var foundUser models.User

	err := c.BindJSON(&user)
	if err != nil {
		return
	}
	res := usersCollection.FindOne(ctx, bson.M{"email": user.Email})
	println(res)
	if err != nil {
		msg := fmt.Sprintf("Error checking email")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "msg": msg})
		return
	}
	if res == nil {
		msg := fmt.Sprintf("No user found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": msg})
		return
	}
	err = res.Decode(&foundUser)

	if err != nil {
		msg := fmt.Sprintf("Invalid email or Password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": msg})
		return

	}
	passwordCheck, msg := VerifyPassword(user.Password, foundUser.Password)
	defer cancel()
	if passwordCheck != true {
		c.JSON(http.StatusUnauthorized, gin.H{"message": msg})
		return
	}
	if foundUser.Email == "" {
		msg := fmt.Sprintf("No user found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": msg})
		return
	}
	token, refreshToken, _ := helpers.GenerateAllToken(foundUser.Email, foundUser.Username, foundUser.UserId, foundUser.UserType)
	helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)
	results := usersCollection.FindOne(ctx, bson.M{"userId": foundUser.UserId})

	err = results.Decode(&foundUser)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, foundUser)
}

//TODO
//Forgot Password
//Verify Email
//Change Password
