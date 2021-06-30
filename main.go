package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wchr-aun/wasted-backend/config"
	"github.com/wchr-aun/wasted-backend/endpoints"
	"github.com/wchr-aun/wasted-backend/middleware"
)

func main() {
	router := gin.Default()
	firebaseAuth := config.SetupFirebase()
	dynamodbCon := config.ConnectDynamoDB()

	router.Use(middleware.CORSMiddleware())
	router.Use(func(c *gin.Context) {
		c.Set("firebaseAuth", firebaseAuth)
	})
	router.Use(func(c *gin.Context) {
		c.Set("dynamodbCon", dynamodbCon)
	})

	// public routes
	// router.GET("", endpoints)

	// private routes
	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware)
	private.GET("/auth", endpoints.GetAuthentication)
	private.POST("/auth", endpoints.PostAuthentication)

	// waste routes
	wasteHandler := new(endpoints.WasteHandler)
	wasteRouter := private.Group("/waste")
	wasteRouter.GET("getMaster", wasteHandler.GetMasterWasteType)
	wasteRouter.GET("getWaste", wasteHandler.GetWasteSeller)

	router.Run()
}
