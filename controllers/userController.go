package controllers

import (
	"context"
	"github.com/bostigger/go-jwt-mongo/database"
	"github.com/bostigger/go-jwt-mongo/helpers"
	"github.com/bostigger/go-jwt-mongo/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

var usersCollection *mongo.Collection = database.GetCollection(database.Client, "users")

func GetUserByID(c *gin.Context) {
	userId := c.Params.ByName("user_id")

	err := helpers.CheckUserAccess(c, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "You cant access this"})
		return
	}
	var user models.User
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Error decoding users"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	res := usersCollection.FindOne(ctx, bson.M{"userId": userId})
	if err != nil {
		return
	}
	err = res.Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Error decoding user"})
		return
	}
	c.JSON(http.StatusOK, user)

}

func GetUsers(c *gin.Context) {
	err := helpers.CheckUserAccess(c, "admin")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You cant access this page"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var users []bson.M

	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, err1 := strconv.Atoi(c.Query("page"))
	if err1 != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(c.Query("startIndex"))

	matchStage := bson.D{{"$match", bson.D{{}}}}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "null"},
			{"total_count", bson.D{
				{"$sum", 1},
			}},
			{"data", bson.D{
				{"$push", "$$ROOT"},
			}},
		}},
	}
	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		}},
	}

	cursor, err := usersCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Error fetching users"})
		return
	}
	err = cursor.All(ctx, &users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Error decoding users"})
		return
	}

	c.JSON(http.StatusOK, users[0])
}
