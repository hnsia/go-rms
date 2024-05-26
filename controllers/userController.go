package controller

import (
	"context"
	"go-rms/database"
	"go-rms/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page-1) * recordPerPage
		startIndex, err2 := strconv.Atoi(c.Query("startIndex"))
		if err2 != nil || startIndex < 1 {
			startIndex = 0
		}

		matchStage := bson.D{{"$match", bson.D{{}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}},
		}

		result, aggregateErr := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage,
		})
		defer cancel()
		if aggregateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing users"})
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred when listing user"})
		}
		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// convert JSON data coming from postman to something that goland understands

		// validate the data based on user struct

		// check if email has already been used by another user

		// hash password

		// check if phone number has already been used by another user

		// get some extra details - created_at, updated_at, ID

		// generate token and refresh token

		// if all ok, insert new user into user collection

		// return status OK and send result back
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// convert login JSON data from postman to golang readable format

		// find user with that email and see if that user exists

		// verify the password

		// if all ok, generate token

		// update tokens - token and refresh token

		// return status OK
	}
}

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	
}