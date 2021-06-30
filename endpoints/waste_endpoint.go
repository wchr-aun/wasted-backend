package endpoints

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gin-gonic/gin"
	"github.com/wchr-aun/wasted-backend/models"
	"log"
	"net/http"
	"os"
)

var MASTER_WASTE_TYPE_TABLE = os.Getenv("MASTER_WASTE_TYPE_TABLE")
var SELLER_WASTE_TABLE = os.Getenv("SELLER_WASTE_TABLE")

type WasteHandler struct{}

func (wh *WasteHandler) GetMasterWasteType(c *gin.Context) {

	dynamodbCon := c.MustGet("dynamodbCon").(*dynamodb.DynamoDB)

	params := &dynamodb.ScanInput{
		TableName: aws.String(MASTER_WASTE_TYPE_TABLE),
	}

	result, err := dynamodbCon.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var wasteTypes []models.MasterWasteType
	for _, i := range result.Items {
		item := models.MasterWasteType{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Printf("Got error unmarshalling: %s", err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		wasteTypes = append(wasteTypes, item)
	}
	c.JSON(200, wasteTypes)
	return
}

func (wh *WasteHandler) GetWasteSeller(c *gin.Context) {
	uuid := c.MustGet("UUID").(string)
	dynamodbCon := c.MustGet("dynamodbCon").(*dynamodb.DynamoDB)

	filter := expression.Name("sellerId").Equal(expression.Value(uuid))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		log.Printf("Got error building expression: %s", err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(SELLER_WASTE_TABLE),
	}

	result, err := dynamodbCon.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var wastes []models.SellerWaste
	for _, i := range result.Items {
		item := models.SellerWaste{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Printf("Got error unmarshalling: %s", err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		wastes = append(wastes, item)
	}
	c.JSON(200, wastes)
	return
}
