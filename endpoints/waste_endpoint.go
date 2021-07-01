package endpoints

import "C"
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
	"strconv"
)

var MASTER_WASTE_TYPE_TABLE = os.Getenv("MASTER_WASTE_TYPE_TABLE")
var SELLER_WASTE_TABLE = os.Getenv("SELLER_WASTE_TABLE")

type WasteHandler struct {
	WasteTypes map[int]models.MasterWasteType
}

func InitWasteHandler(dynamodbCon *dynamodb.DynamoDB) WasteHandler {
	wh := WasteHandler{}
	wh.WasteTypes = wh.getMasterWasteType(dynamodbCon)
	return wh
}

func (wh *WasteHandler) getMasterWasteType(dynamodbCon *dynamodb.DynamoDB) map[int]models.MasterWasteType {

	params := &dynamodb.ScanInput{
		TableName: aws.String(MASTER_WASTE_TYPE_TABLE),
	}

	result, err := dynamodbCon.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		return nil
	}

	wasteTypes := make(map[int]models.MasterWasteType)
	for _, i := range result.Items {
		item := models.MasterWasteType{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Printf("Got error unmarshalling: %s", err)
			return nil
		}
		wasteTypes[item.WasteId] = item
	}

	return wasteTypes
}

func (wh *WasteHandler) GetMasterWasteType(c *gin.Context) {
	c.JSON(200, wh.WasteTypes)
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

func (wh *WasteHandler) UpdateWasteSeller(c *gin.Context) {
	uuid := c.MustGet("UUID").(string)
	dynamodbCon := c.MustGet("dynamodbCon").(*dynamodb.DynamoDB)

	waste := models.SellerWaste{SellerId: uuid}
	err := c.ShouldBind(&waste)
	if err != nil {
		log.Printf("Got error binding waste : %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	_, found := wh.WasteTypes[waste.WasteId]
	if !found {
		log.Printf("Got error matching waste in master waste type")
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Not found waste type in master"})
		return
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				N: aws.String(strconv.Itoa(waste.Amount)),
			},
		},
		TableName: aws.String(SELLER_WASTE_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"wasteId": {
				N: aws.String(strconv.Itoa(waste.WasteId)),
			},
			"sellerId": {
				S: aws.String(uuid),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set amount = :r"),
	}

	result, err := dynamodbCon.UpdateItem(input)
	if err != nil {
		log.Printf("Got error calling UpdateItem: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, &waste)
	if err != nil {
		log.Printf("Got error Unmarshal: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, waste)
}
