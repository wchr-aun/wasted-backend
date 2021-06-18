package endpoints

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/wchr-aun/wasted-backend/models"
)

func GetAuthentication(c *gin.Context) {
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
		log.Printf("Got error calling GetItem: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	if result.Item == nil {
		err := registerUser(c, models.UserTable{Uuid: uuid})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, models.UserResponse{
			New: true,
		})
		return
	}

	userResponse := models.UserResponse{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &userResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusOK, userResponse)
}

func PostAuthentication(c *gin.Context) {
	uuid := c.MustGet("UUID").(string)

	addingUser := models.UserTable{
		Uuid: uuid,
	}
	c.ShouldBind(&addingUser)

	registerUser(c, addingUser)

	userResponse := models.UserResponse{
		Fullname:    addingUser.Fullname,
		PhoneNumber: addingUser.PhoneNumber,
		PhotoUrl:    addingUser.PhotoUrl,
	}

	c.JSON(http.StatusOK, userResponse)
}

func registerUser(c *gin.Context, user models.UserTable) gin.H {
	dynamodbCon := c.MustGet(("dynamodbCon")).(*dynamodb.DynamoDB)

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Printf("Got error marshalling new movie item: %s", err)
		return gin.H{"title": "Internal Server Error", "msg": "Please contact admins for the help"}
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("users"),
	}

	_, err = dynamodbCon.PutItem(input)
	if err != nil {
		log.Printf("Got error calling PutItem: %s", err)
		return gin.H{"title": "Internal Server Error", "msg": "Please contact admins for the help"}
	}

	return nil
}
