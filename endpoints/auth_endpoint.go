package endpoints

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/wchr-aun/wasted-backend/models"
)

var USERS_TABLE string

func init() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	USERS_TABLE = os.Getenv("USERS_TABLE")
}

func Authentication(c *gin.Context) {
	uuid := c.MustGet("UUID").(string)
	dynamodbCon := c.MustGet(("dynamodbCon")).(*dynamodb.DynamoDB)

	result, err := dynamodbCon.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USERS_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {
				S: aws.String(uuid),
			},
		},
	})

	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		addUser(c)
		return
	}

	user := models.User{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusOK, user)
}

func addUser(c *gin.Context) {
	uuid := c.MustGet(USERS_TABLE).(string)
	dynamodbCon := c.MustGet(("dynamodbCon")).(*dynamodb.DynamoDB)

	user := models.User{
		Uuid:        uuid,
		Fullname:    "test user",
		PhoneNumber: "0123456789",
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("users"),
	}

	_, err = dynamodbCon.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	user.New = true

	c.JSON(http.StatusOK, user)
}
